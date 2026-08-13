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
- **Job model**: `event`/`title` (NOT `event_id`/`event_name`), `Elapsed` is float, `code` may be string
- **APIKey create**: returns `plain_key` (only exposure); surfaced on returned `APIKey.Key`
- **get_multiple**: `servers` is an OBJECT keyed by ID; lists also include `channels`, `buckets`, `secrets`; flags `alerts`/`stats`/`serverCache`
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
- **Key resolution**: The provider passes human-readable keys through to SDK methods; `sdk.Resolver` exists for key→ID lookups but lives in the SDK layer and is NOT used by the provider. Do NOT resolve keys via plan modifiers.
- **MCP library**: Use `github.com/modelcontextprotocol/go-sdk` (official, not mark3labs)
- **go.mod**: `github.com/gumieri/xyops-integrations`

## CI & Lint

- **Linting (mandatory)**: Run `golangci-lint run ./...` after any change. CI blocks PRs on lint findings. Config lives in `.golangci.yml` (v2 schema): `gocyclo` ≥15, `funlen` 150/80, `nestif` 4, `staticcheck`, `govet` (all analyzers), `errcheck`, `misspell`, plus `gofmt`/`goimports` formatters. Test files (`*_test.go`) are excluded from linting.
- **CI**: `.github/workflows/ci.yml` runs `go test ./... -race`, `go vet ./...`, and golangci-lint on every push to `main` and every PR.
- **Security scanning**: `.github/workflows/codeql.yml` runs CodeQL for Go on push/PR to `main` plus a weekly schedule (Mon 04:00 UTC).
- **Dependencies**: `.github/dependabot.yml` opens weekly PRs for Go modules and GitHub Actions.
- **Release**: Pushing a `vX.Y.Z` tag triggers `.github/workflows/release.yml`, which builds multi-arch container images (GHCR), runs GoReleaser (binaries, deb/rpm/arch, AUR), generates SPDX SBOMs, and cosign-signs images + checksums. GoReleaser itself runs `go test -race` and `go vet` before publishing.

## Build & Test

```bash
go build ./...
go test ./... -count=1
go vet ./...
golangci-lint run ./...
```

Mise tasks: `mise run build|test|race|lint|vet|validate-all|snapshot|release`

Acceptance tests: `TF_ACC=1 go test ./internal/provider/... -v -count=1 -run "^TestAcc"` (requires `XYOPS_BASE_URL` + `XYOPS_API_KEY`)