# xyops-integrations — Project Context

## Project Structure

This Go module (`github.com/gumieri/xyops-integrations`) builds three things from the same SDK:

| Directory | Purpose |
|-----------|---------|
| `sdk/` | Go SDK — typed API client for xyops conductor REST API |
| `internal/provider/` | Terraform provider (terraform-plugin-framework) |
| `internal/mcp/` | MCP server (official MCP Go SDK) |
| `cmd/terraform-provider-xyops/` | Provider binary entry point |
| `cmd/xyops-mcp/` | MCP binary entry point |

## Architecture Rule

The SDK (`sdk/`) absorbs ALL API quirks so callers get clean Go methods.
Both the TF provider and MCP server use the same `sdk.Client` — the SDK is the only layer that talks to the HTTP API.

## Implementation Order

1. Phase 011: SDK — Models + HTTP Client (no external deps beyond stdlib)
2. Phase 012: SDK — All resource CRUD + List + Job Services
3. Phase 013: TF Provider — Scaffold + Tag resource (acceptance test as primary)
4. Phase 014: TF Provider — All remaining resources + data sources
5. Phase 015: MCP Server — All tools
6. Phase 016: Integration tests + docs + GoReleaser

## Key Gotchas

- **Error envelope**: `code` (int 0=ok OR string like `"session"`/`"test"`) + `description` (NOT `message`)
- **Session detection**: `code:"session"` (string) on HTTP 200 = auth error; `APIError.IsAuth()` checks code string + 401/403
- **search_jobs params**: `?query=`, `?sort_by=`, `?sort_dir=` — NOT `q`/`sort`
- **get_active_jobs filter**: `event`, NOT `event_id`
- **get_workflow_job_summary filter**: `workflow.job`, NOT `id`
- **get_jobs response**: `jobs` array, NOT `rows`
- **get_servers**: GET `get_servers/v1` (not POST); returns `servers` OBJECT keyed by ID + `masters`; `ServersService.List` flattens to sorted `[]Server`
- **Server model**: `info.os/cpu/memory/virt` nested; NOT flat `os`/`arch`/`status` (legacy `Satellite` type kept for compat)
- **satellite/config**: plain JSON config.json document, NO `{code, data}` envelope
- **Health**: `ping/v1`, NOT `/health`
- **Job model**: `event`/`title` (NOT `event_id`/`event_name`), `Elapsed` is float, `Code` may be string (`0` or `"warning"`/`"critical"`/`"abort"`), `Created`/`Completed`/`Started` are **fractional** epoch (float64)
- **APIKey create**: returns `plain_key` (only exposure); surfaced on returned `APIKey.Key`
- **get_multiple**: **Restored in v1.0.90** (removed in v1.0.49). `Epoch` is fractional float64 in v1.0.90 (was int64). SDK handles both.
- **Privileges**: `sdk.BoolMap` (custom type) — JSON values may be `true`/`false` OR `1`/`0`. SDK unmarshals both to `map[string]bool` transparently.
- **APIKey.Active**: `int` (1/0), NOT bool
- **WebHook.Headers**: `[]Header{{Name, Value}}`, NOT `map[string]string`
- **APIKey.Key**: write-only, omit on update payload
- **APIKey read**: POST body `{"id":"..."}`, NOT GET
- **Plugin upsert**: NO upsert — always `create_plugin/v1` even with `id`
- **Plugin.Enabled**: `sdk.EnabledInt` — JSON may be `true`/`false` OR `1`/`0`. SDK unmarshals both to int.
- **PluginParam.Value**: polymorphic (`any`) — may be string, number, or bool. Provider renders as string.
- **PluginParam.Required** / **PluginParam.Locked**: nullable `*bool`. Provider treats nil as false.
- **Privilege errors**: v1.0.90 returns `code:"access"` (HTTP 403) for missing privilege. SDK treats it as a plain `APIError`; `IsAuth()` is false (it's a permission denial, not an invalid session). Cross-version: all polymorphic types (`BoolMap`, `EnabledInt`, `*bool`, fractional epochs) verified against BOTH v1.0.49 and v1.0.90.
- **Update response**: `write_returns_object = false` — must Get after Update
- **Monitor read**: GET `get_monitor/v1?id=X` (singular exists, no list+filter needed)
- **Monitors refactored**: No longer `get_monitors/v1` POST list+filter. Use `get_monitor/v1?id=X` directly.
- **get_job/v1 shape**: Returns a workflow-graph document (`{code, job:{type,workflow:{nodes,connections,state,jobs}}}`) for workflow jobs, NOT the flat `Job` record used by `search_jobs`/`get_jobs`. `JobsService.Get` unmarshals into `Job` and will yield mostly-empty fields for workflows. Use `search_jobs`/`get_jobs` for flat job records.
- **Limits/Actions (empirical, v1.0.49)**: `Event.Limits[].type` values observed in production: `job`, `queue`, `retry`, `time`. `Event.Actions[].condition` observed: `success` (docs mention `error`/`warning`; don't assume only one).
- **List endpoints**: Exist for ALL resource types via GET `get_*/v1`
- **Key resolution**: The provider passes human-readable keys through to SDK methods; `sdk.Resolver` exists for key→ID lookups but lives in the SDK layer and is NOT used by the provider. Do NOT resolve keys via plan modifiers.
- **MCP library**: Use `github.com/modelcontextprotocol/go-sdk` (official, not mark3labs)
- **go.mod**: `github.com/gumieri/xyops-integrations`

## Documentation & Privacy

This is a public repository — assume anything committed is permanently public
(history is not rewritten by default).

- **No private infrastructure in tracked files**: never write internal
  hostnames, URLs, endpoints, workspace slugs, tool names, or credentials into
  code, tests, docs, commit messages, or PR text. Use the `nxample.com`
  placeholder convention for all endpoints (see CONTRIBUTING.md).
- **Internal design docs live outside this repo** (maintainer's private
  knowledge base, linked from the issue tracker). `docs/` is
  operator/user-facing reference only — design deliberation, decision
  history, and withdrawn plans do not belong here.
- **Ignored dirs stay untracked**: `.opencode/`, `.scratch-mcp/`,
  `.scratch-opencode/`, `.env*`, `*.pem`, `*.key` hold local agent config and
  live credentials. Ignore rules do not untrack already-added files, and
  `git add -f` bypasses them — check `git status` / `git ls-files` before
  committing any newly created config or scratch file.

## Claim Discipline

- **Verify before recording**: confirm claims about the codebase (grep,
  tests, `git log`) before writing them into docs, CHANGELOG entries, PR
  descriptions, or issue comments. Unverified "shipped"/"implemented" claims
  are the top historical failure mode.
- **Working tree ≠ git index**: filesystem tools (grep/ls) see the tree;
  they say nothing about what git tracks. Tracking questions need
  `git ls-files` or the hosting API.
- **Timing flakes**: perf/benchmark-based tests are noisy on shared CI
  runners — rerun and confirm the branch isn't stale against `main` before
  treating a single failure as a regression.

## CI & Lint

- **Linting (mandatory)**: Run `golangci-lint run ./...` after any change. CI blocks PRs on lint findings. Config lives in `.golangci.yml` (v2 schema): `gocyclo` ≥15, `funlen` 150/80, `nestif` 4, `staticcheck`, `govet` (all analyzers), `errcheck`, `misspell`, plus `gofmt`/`goimports` formatters. Test files (`*_test.go`) are excluded from linting.
- **CI**: `.github/workflows/ci.yml` runs `go test ./... -race`, `go vet ./...`, and golangci-lint on every push to `main` and every PR.
- **Security scanning**: `.github/workflows/codeql.yml` runs CodeQL for Go on push/PR to `main` plus a weekly schedule (Mon 04:00 UTC).
- **Dependencies**: `.github/dependabot.yml` opens weekly PRs for Go modules and GitHub Actions.
- **Release**: Pushing a `vX.Y.Z` tag triggers `.github/workflows/release.yml`, which builds multi-arch container images (GHCR), runs GoReleaser (binaries, deb/rpm/arch, AUR), generates SPDX SBOMs, and cosign-signs images + checksums. GoReleaser itself runs `go test -race` and `go vet` before publishing.
- **Merging workflow PRs**: a `gh` CLI token without the `workflow` scope cannot merge or update PRs that touch `.github/workflows/*` — apply such changes via a direct push and close the PR as superseded (note the reason in a comment).

## Build & Test

```bash
go build ./...
go test ./... -count=1
go vet ./...
golangci-lint run ./...
```

Mise tasks: `mise run build|test|race|lint|vet|validate-all|snapshot|release`

Acceptance tests: `TF_ACC=1 go test ./internal/provider/... -v -count=1 -run "^TestAcc"` (requires `XYOPS_BASE_URL` + `XYOPS_API_KEY`)

## Testing Policy (XYOINT-15)

xyops ships no machine-readable API spec, so external truth is **captured**, never hand-duplicated:

- **Contract split**: Go tests assert the SDK contract only (wire round-trip, envelope parsing, quirk absorption). xyops *behavior* (required params, validation policy) is smoke-tested and documented in `docs/api-reference.md` — behavioral assertions in Go tests are duplicity that masks drift.
- **Golden wire captures**: `XYOPS_CAPTURE_GOLDENS=1` re-fetches read-only endpoints into `test/testdata/goldens/<version>/`; CI (`TestGoldenWireContracts`, offline) asserts the SDK still unmarshals every captured body. Re-capture on every xyops version bump — the capture diff is the API changelog.
- **Version pinning**: `TestConductorVersion` logs the live conductor version; `XYOPS_EXPECT_VERSION=<ver>` fails on mismatch.
- **Dynamic discovery**: tests derive instance data (categories, plugins, groups) from the instance itself — never hardcode names or IDs.