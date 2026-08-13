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

- **Privileges**: `map[string]int` (values 1/0), NOT bool
- **APIKey.Active**: `int` (1/0), NOT bool
- **WebHook.Headers**: `[]Header{{Name, Value}}`, NOT `map[string]string`
- **APIKey.Key**: write-only, omit on update payload
- **APIKey read**: POST body `{"id":"..."}`, NOT GET
- **Plugin upsert**: NO upsert — always `create_plugin/v1` even with `id`
- **Update response**: `write_returns_object = false` — must Get after Update
- **Monitor read**: GET `get_monitor/v1?id=X` (singular exists, no list+filter needed)
- **Monitors refactored**: No longer `get_monitors/v1` POST list+filter. Use `get_monitor/v1?id=X` directly.
- **List endpoints**: Exist for ALL resource types via GET `get_*/v1`
- **Key resolution**: At apply time via SDK resolver + list endpoints. NOT via plan modifiers.
- **MCP library**: Use `github.com/modelcontextprotocol/go-sdk` (official, not mark3labs)
- **go.mod**: `github.com/gumieri/xyops-integrations`

## Build & Test

```bash
go build ./...
go test ./... -count=1
go vet ./...
```

Acceptance tests: `TF_ACC=1 go test ./internal/provider/... -v -count=1 -run "^TestAcc"`