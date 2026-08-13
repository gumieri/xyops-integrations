# Contributing to xyops-integrations

## Project Structure

```
sdk/                 # Go SDK — typed API client for the xyops conductor REST API
internal/provider/   # Terraform provider (terraform-plugin-framework v1.x)
internal/mcp/        # MCP server (official modelcontextprotocol/go-sdk)
cmd/terraform-provider-xyops/  # Provider binary entry point
cmd/xyops-mcp/                  # MCP binary entry point
docs/                # Documentation (API reference)
examples/            # Usage examples (examples/mcp, examples/terraform)
test/                # Integration tests (require live xyops instance + env vars)
```

## Architecture Rule

The SDK (`sdk/`) absorbs ALL API quirks so callers get clean Go methods. Both the TF provider and MCP server use the same `sdk.Client` — the SDK is the only layer that talks to the HTTP API.

## Commit Conventions

Use [Conventional Commits 1.0.0](https://www.conventionalcommits.org/):

```
feat(sdk): add support for job search pagination
fix(provider): resolve group keys on create
docs: add API reference
test: cover plugin no-upsert behavior
```

Common types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `build`, `ci`, `chore`.

## Adding a New Resource

Follow this order — each layer depends on the one before it:

1. **SDK models** (`sdk/models.go`) — add the typed struct with correct JSON tags.
   Remember the type quirks:
   - `Privileges`: `map[string]int` (values 1/0), NOT bool
   - `APIKey.Active`: `int` (1/0), NOT bool
   - `WebHook.Headers`: `[]Header{{Name, Value}}`, NOT `map[string]string`
   - `APIKey.Key`: write-only, omit on update payload
2. **SDK service** (`sdk/<resource>.go`) — add `List`/`Get`/`Create`/`Update`/`Delete`.
   - `Update` must call `Get` afterward (`write_returns_object = false`).
   - Plugins have NO upsert — always `create_plugin/v1`.
   - APIKey read uses POST body `{"id":"..."}`, not GET.
3. **SDK tests** (`sdk/<resource>_test.go`) — mock the HTTP server with `setupTestServer`; verify each method's request path and response parsing.
4. **Terraform resource** (`internal/provider/resource_<name>.go`) — implement the `resource.Resource` interface. Use plan modifiers where appropriate (e.g. `stringplanmodifier.UseStateForUnknown()` for computed fields; `RequiresReplace` for immutable fields). The provider passes keys through to the SDK methods; the `sdk.Resolver` type exists for key→ID lookups and is kept in the SDK layer.
5. **Terraform data source** (if read-only access is useful) — `datasource_<name>.go`.
6. **Provider registration** — add to `Resources()`/`DataSources()` in `internal/provider/provider.go`.
7. **Acceptance tests** — `internal/provider/resource_<name>_test.go` using `ProtoV6ProviderFactories`. Run with `TF_ACC=1` and `XYOPS_BASE_URL` + `XYOPS_API_KEY` set.
8. **MCP tools** (optional) — add handlers in `internal/mcp/tools_<name>.go` using the typed `adaptResult` pattern.

## Running Tests

```bash
# Unit tests (no external deps)
go test ./... -count=1

# Race detector
go test -race -count=1 ./...

# Lint and vet
golangci-lint run ./...
go vet ./...

# Acceptance tests (require TF_ACC=1, live xyops instance, and XYOPS_BASE_URL + XYOPS_API_KEY)
TF_ACC=1 XYOPS_BASE_URL=https://xyops.example.com XYOPS_API_KEY=... go test ./internal/provider/... -v -count=1 -run "^TestAcc"

# Integration tests (require XYOPS_BASE_URL + XYOPS_API_KEY)
XYOPS_BASE_URL=https://xyops.example.com XYOPS_API_KEY=... go test ./test/... -count=1
```

## Development Setup

```bash
# Lint (uses golangci-lint; config is optional, defaults apply)
golangci-lint run ./...

# Build all binaries
mise run build

# Snapshot release (builds to dist/)
goreleaser release --snapshot --clean

# Full release (requires GitHub token)
goreleaser release --clean
```

## Environment Variables

| Variable | Purpose | Required |
|----------|---------|----------|
| `XYOPS_BASE_URL` | xyops conductor base URL (must include an `http`/`https` scheme and host, with no query/fragment; required by the MCP binary, provider, and acceptance/integration tests) | Yes |
| `XYOPS_API_KEY` | API key for authentication | Yes |
| `TF_ACC` | Enables Terraform acceptance tests | Acceptance tests |

## Reviewer Guidance

- The SDK must remain the only layer that talks to the API.
- Never add secrets, keys, or real hostnames to code, tests, or docs — use `xyops.example.com` placeholders.
- Run `gofmt -w .` before committing.
- Keep commits focused: one logical change per commit.
