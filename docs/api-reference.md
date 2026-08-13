# xyops API Reference

This document catalogs the xyops conductor REST API as consumed by the Go SDK in `sdk/`. The SDK absorbs all API quirks so callers get clean Go methods; this reference documents the underlying wire behavior.

All endpoints are under `/api/app/`.

- **`code`**: `0` = success, non-zero = error. The string `"session"` (or `"test"`) is also a valid error code — the SDK accepts both numeric and string codes.
- **HTTP status**: non-2xx responses are treated as errors; the response body is parsed for `code`/`description` when present.
- **Auth**: API keys are passed via the `X-API-Key` header — the SDK always sends it, there is no override. Keys are typically 24-character alphanumeric and case-sensitive (API requires a minimum of 16 characters).

## CRUD Resources

| Resource | List | Get | Create | Update | Delete | Quirks |
|----------|------|-----|--------|--------|--------|--------|
| **Tag** | `get_tags/v1` (GET) | `get_tag/v1?id=X` (GET) | `create_tag/v1` (POST) | `update_tag/v1` (POST) | `delete_tag/v1` (POST) | ID prefix `t...` |
| **Group** | `get_groups/v1` (GET) | `get_group/v1?id=X` (GET) | `create_group/v1` (POST) | `update_group/v1` (POST) | `delete_group/v1` (POST) | ID prefix `g...` |
| **Category** | `get_categories/v1` (GET) | `get_category/v1?id=X` (GET) | `create_category/v1` (POST) | `update_category/v1` (POST) | `delete_category/v1` (POST) | ID prefix `ca...` |
| **Plugin** | `get_plugins/v1` (GET) | `get_plugin/v1?id=X` (GET) | `create_plugin/v1` (POST) | `update_plugin/v1` (POST) | `delete_plugin/v1` (POST) | ID prefix `p...`; **NO upsert** — always `create_plugin/v1` even with `id` |
| **Event** | `get_events/v1` (GET) | `get_event/v1?id=X` (GET) | `create_event/v1` (POST) | `update_event/v1` (POST) | `delete_event/v1` (POST) | ID prefix `e...`; workflows: `plugin="_workflow"`, `type="workflow"`, nested `workflow` DAG |
| **Monitor** | `get_monitors/v1` (GET) | `get_monitor/v1?id=X` (GET) | `create_monitor/v1` (POST) | `update_monitor/v1` (POST) | `delete_monitor/v1` (POST) | ID prefix `m...` |
| **Alert** | `get_alerts/v1` (GET) | `get_alert/v1?id=X` (GET) | `create_alert/v1` (POST) | `update_alert/v1` (POST) | `delete_alert/v1` (POST) | ID prefix `a...` |
| **Role** | `get_roles/v1` (GET) | `get_role/v1?id=X` (GET) | `create_role/v1` (POST) | `update_role/v1` (POST) | `delete_role/v1` (POST) | `privileges` = `map[string]int` (1/0, NOT bool) |
| **API Key** | `get_api_keys/v1` (GET) | `get_api_key/v1` (**POST** body `{"id":"..."}`) | `create_api_key/v1` (POST) | `update_api_key/v1` (POST) | `delete_api_key/v1` (POST) | **QUIRKS**: (1) read is POST, not GET; (2) `key` write-only, NOT sent on update; (3) `active` is int `1`/`0` |
| **Web Hook** | `get_web_hooks/v1` (GET) | `get_web_hook/v1?id=X` (GET) | `create_web_hook/v1` (POST) | `update_web_hook/v1` (POST) | `delete_web_hook/v1` (POST) | `headers` = list of `{name, value}` objects, NOT a map |

ID prefixes (e.g. `t...`, `g...`) are API conventions observed in practice; the SDK does NOT validate or enforce them.

**Upsert pattern**: All resources except plugins route to `update_*/v1` when the caller supplies an `id`, otherwise `create_*/v1`. This enables adopting existing resources without import. Plugins ALWAYS use `create_plugin/v1` regardless of `id`.

**Update response**: `write_returns_object = false` — update responses do NOT include the object. The SDK calls `Get` after every update and returns the fresh object.

## Config (singleton)

| Operation | Endpoint | Notes |
|-----------|----------|-------|
| Get | `admin_get_config/v1` (POST, no body) | Returns full config map |
| Update | `admin_update_config/v1` (POST) | Sparse merge; send empty arrays to clear fields |

No delete — clear fields via update. In Terraform this maps to `prevent_destroy`.

## Servers & Satellite

| Operation | Endpoint | Notes |
|-----------|----------|-------|
| List servers | `get_servers/v1` (GET) | Returns `servers` **OBJECT** keyed by server ID (NOT an array) plus `masters`. `ServersService.List` flattens to a sorted `[]Server`. |
| List masters | `get_servers/v1` (GET) | `masters` object keyed by host ID with `Conductor` status/version/stats — via `ServersService.Masters`. |
| Satellite config | `satellite/config?t=X` (GET) | Returns the raw `config.json` document (`port`, `secure`, `hosts`, `server_id`, `auth_token`, `initial`) — NO `{code,data}` envelope. |
| Satellite token | `get_satellite_token/v1` (POST) | Generate/refresh bootstrap token; returns under `data` key |
| Satellite releases | `get_satellite_releases/v1` (GET) | List available binary releases |

## Job & Trigger Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `run_event/v1` | POST | **JOB DISPATCH** — trigger event execution on demand. Supports `params`, `input`, `tags`, `test: true`. Match by `id` or `title`. Returns `{"code":0,"id":"j..."}` where `id` is the job ID; call `get_job/v1?id=j...` for details and token. |
| `get_active_jobs/v1` | GET | List active jobs (queued/ready/active/finishing) with filters + pagination. Filter by `event` (NOT `event_id`); also `state`, `server`, `workflow.job`. |
| `get_active_job_summary/v1` | GET | Summarize active jobs by event/state/source/targets |
| `get_workflow_job_summary/v1` | GET | Summarize workflow jobs by node (`workflow.job` filter param required) |
| `get_job/v1` | GET | Single job details (`id` param; returns `job` + `token` for log access). For **workflow** jobs the `job` document is a graph (`{type, workflow:{nodes,connections,state,jobs}}`), not a flat `Job` record — use `search_jobs`/`get_jobs` for flat records. |
| `get_jobs/v1` | POST | Batch fetch jobs by IDs (`ids` array, optional `verbose: true`). Response uses `jobs` array (NOT `rows`). |
| `get_event_history/v1` | GET | Revision history for an event (paginated) |
| `search_jobs/v1` | GET | Search completed jobs (PxQL/Unbase-style query syntax). Params: `query` (NOT `q`), `offset`, `limit`, `sort_by` (NOT `sort`), `sort_dir` (`1`/`-1`), `verbose`, `select`. |
| `abort_job/v1` | POST | Abort running job |
| `resume_job/v1` | POST | Resume suspended job (optional `params`) |
| `delete_job/v1` | POST | Delete completed job (logs + files) |
| `flush_event_queue/v1` | POST | Flush all queued jobs for an event |
| `manage_job_tags/v1` | POST | Replace tags on a completed job |

## Batch Endpoint

> **Note (v1.0.49)**: `get_multiple/v1` is **no longer available** — the conductor returns `"Unsupported API: app/api_get_multiple"`. This section is retained for historical/backwards-compatibility reference only. Prefer individual list endpoints (`get_tags/v1`, `get_groups/v1`, ...) which all work on v1.0.49.

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `get_multiple/v1` | GET | Batch-fetch several list types in one call: `?lists=events,categories,plugins&jobs=1&alerts=1&state=1&stats=1&servers=1&serverCache=1`. Lists: `events`, `categories`, `plugins`, `tags`, `groups`, `monitors`, `alerts`, `roles`, `api_keys`, `web_hooks`, `channels`, `buckets`, `secrets`, or `all`. `servers` is an **OBJECT** keyed by ID; lists are arrays; `epoch` is always present. |

## Health

`GET /api/app/ping/v1` — returns `{"code": 0}` when reachable. The SDK considers the instance healthy if `code == 0`. HTTP 503 indicates reachable but not the primary instance.

## Response Structure

```json
{"code": 0, "tag": {"id": "tmnditeue2pu40sc", "title": "backup", ...}}
```

- Response key varies per resource: `tag`, `group`, `event`, `monitor`, `alert`, `role`, `api_key`, `web_hook`, `category`, `plugin`, `config`, `servers`, and `data` (used by `get_satellite_token/v1`).
- Create responses include the full object (`create_returns_object = true`).
- List responses use `rows` (array) and `list.length` (total count pre-pagination).

## Error Handling

- `code != 0` or a string code (e.g. `"session"`, `"test"`) → structured error with the `description` field (a legacy `message` field is also honored as a fallback).
- HTTP 4xx/5xx → JSON body parsed for `code`/`description`; fallback uses raw body.
- Auth errors: HTTP 401/403, or a string code of `"session"` (e.g. on a 200 response) — `APIError.IsAuth()` covers all three.
- The SDK exposes `sdk.APIError{Code (any: int|string), Description, HTTPStatus}` with `IsAuth()`.

## Field Reference by Resource

**Tag**: `id`, `title`

**Group**: `id`, `title`, `hostname_match`, `sort_order`

**Category**: `id`, `title`, `enabled`

**Plugin**: `id`, `title`, `type` (`"event"|"monitor"|"action"|"scheduler"`), `command`, `script`, `enabled`, `params`, `uid`, `gid`, `groups`, `format`

**Event**: `id`, `title`, `plugin` (ID ref), `targets` (group ID list), `category` (ID ref), `algo`, `enabled`, `fields`, `triggers`, `limits`, `actions`, `tags` (tag ID list), `notes`, `params`, `expression`, `icon`, `type` (`"event"|"workflow"`)

**Default limits** (SDK injects on create when empty — not on update): `[{"type":"time","enabled":true,"duration":1800}]` (per the API `search_jobs` example; old SDK used `job`/`queue`/`retry`)

**Limit types**: `"time"`, `"job"`, `"queue"`, `"retry"` (API example uses `"time"`)

**Action conditions**: `"error"` (API example uses `"error"`), also `"success"`, `"always"`

**Trigger types**: `"manual"`, `"schedule"` (calendar arrays + timezone), `"interval"` (duration seconds), `"single"` (epoch unix timestamp)

**Monitor**: `id`, `title`, `source`, `data_type` (default `"integer"`), `display` (default true), `groups`, `suffix`, `min_vert_scale`, `data_match`, `delta`, `divide_time`, `zero_min`, `notes`

**Alert**: `id`, `title`, `enabled`, `expression`, `message` (default `""`), `groups`, `monitor_id`, `max_count` (default 1), `notes`

**Role**: `id`, `title`, `enabled`, `privileges` (**`map[string]int`** — 1/0), `categories`, `notes`

**API Key**: `id`, `title`, `description` (default `""`), `key` (write-only; min 16 chars is an API requirement, NOT validated by the SDK — the MCP binary enforces it at startup; ensure caller compliance; NOT sent on update), `active` (**`int`** — 1/0), `privileges` (**`map[string]int`**), `roles`, `expires`. **Create only**: response includes `plain_key` (the actual key value, ONLY available at create time); SDK surfaces this on `APIKey.Key` for the returned object.

**Web Hook**: `id`, `title`, `enabled`, `url`, `method`, `headers` (**list of `{name, value}` objs**), `body`, `timeout`, `retries`, `follow`, `ssl_cert_bypass`, `max_per_day` (int, 0 = unlimited)

**Server**: `id`, `hostname`, `ip`, `enabled`, `groups`, `title`, `icon`, `autoGroup`, `created`, `modified`, `socket_id`, `keywords`, and nested `info` (`os`{`platform`,`distro`,`release`,`arch`}, `cpu`{`cores`,`combo`}, `memory`{`total`}, `virt`{`vendor`}, `satellite`). `get_servers/v1` returns `servers` as an object keyed by ID.

**Job**: `id`, `event` (NOT `event_id`), `title` (NOT `event_name`), `state`, `target`, `params`, `created`, `completed`, `elapsed` (**float**), `code` (int or string like `"warning"`/`"critical"`/`"abort"`), `description`, `plugin`, `type` (`"event"|"workflow"`), `category`, `source`, `progress`, `tag`, `token`, `server`, `username`, `started`. `get_jobs/v1` returns the batch under a `jobs` array.

## Workflow DAG Model

```json
{
  "nodes": [
    {"id": "n00", "type": "event", "x": 50, "y": 80, "data": {"event": "...", "targets": ["group-id"]}},
    {"id": "t00", "type": "trigger", "x": 50, "y": 20},
    {"id": "mx00", "type": "controller", "x": 200, "y": 50, "data": {"controller": "multiplex", "stagger": 60}}
  ],
  "connections": [
    {"id": "ct00n00", "source": "t00", "dest": "n00", "condition": "success"}
  ]
}
```

Node types: `"event"`, `"trigger"`, `"controller"`. IDs must be strictly alphanumeric — an API requirement. The SDK does NOT validate this; ensure caller compliance or the API will reject the payload.

## Key → ID Resolution

Events, plugins, monitors, alerts, and roles reference other resources by human-readable key. The SDK provides a `Resolver` (`sdk.Resolver`) that:

1. Checks a pre-seeded cache first.
2. On miss, queries the corresponding list endpoint and rebuilds the cache for that type.
3. Returns the internal ID for the key.

Supported resolver types: `tag`, `group`, `category`, `plugin`, `monitor`, `role`, `alert`, `api_key`, `webhook`, `event`.
