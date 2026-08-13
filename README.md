# xyops-integrations

Go SDK, Terraform provider, and MCP server for xyops.

## Overview

This project provides three components:

1. **Go SDK** (`sdk/`) — Typed API client for the xyops conductor REST API. Absorbs all API quirks so callers get clean Go methods.
2. **Terraform provider** (`internal/provider/`) — Native `terraform-plugin-framework` v1.x provider for managing xyops resources.
3. **MCP server** (`internal/mcp/`) — stdio-based MCP server enabling Claude/opencode to query and manage xyops resources via tools.

## Quick Start

### Building

```bash
# Build all binaries
mise run build

# Build individual binaries
mise run build-provider  # terraform-provider-xyops
mise run build-mcp       # xyops-mcp
```

### Terraform Provider

The provider works via `dev_overrides` (local development):

```hcl
# ~/.terraformrc
provider_installation {
  dev_overrides {
    "gumieri/xyops" = "/path/to/xyops-integrations"
  }
}
```

Example usage:

```hcl
# provider.tf
terraform {
  required_providers {
    xyops = {
      source = "gumieri/xyops"
    }
  }
}

provider "xyops" {
  base_url = "https://xyops.example.com"
  api_key  = var.xyops_api_key
}

resource "xyops_tag" "example" {
  title = "backup"
}
```

Import existing resources:

```bash
terraform import xyops_tag.example txxxxxxxx
```

### MCP Server

Configure opencode:

```json
{
  "mcp": {
    "servers": {
      "xyops": {
        "command": "./xyops-mcp",
        "env": {
          "XYOPS_BASE_URL": "https://xyops.example.com",
          "XYOPS_API_KEY": "${XYOPS_API_KEY}"
        }
      }
    }
  }
}
```

Available MCP tools include: `xyops_list_servers`, `xyops_list_events`, `xyops_search_jobs`, `xyops_trigger_event`, `xyops_create_tag`, and many more.

## Architecture

```
xyops-integrations/
├── sdk/                 # Go SDK (typed API client, importable: github.com/gumieri/xyops-integrations/sdk)
├── internal/provider/   # Terraform provider logic
├── internal/mcp/        # MCP server logic
├── cmd/terraform-provider-xyops/main.go  # Provider binary
├── cmd/xyops-mcp/main.go                # MCP binary
```

The SDK is the foundation — both the TF provider and MCP server share the same `sdk.Client`.

## Testing

```bash
# Unit tests (no external deps)
mise run test

# Race detector
mise run race

# Lint and vet
mise run lint
mise run vet

# All validation
mise run validate-all
```

Acceptance tests (require `TF_ACC=1` and live xyops instance):

```bash
TF_ACC=1 go test ./internal/provider/... -v -count=1 -run "^TestAcc"
```

## Development

- Go: 1.26+ (via mise)
- Lint: `golangci-lint run ./...`
- Format: `gofmt -w .`
- Release: `goreleaser release --clean` (or `--snapshot` for testing)

## Registry Publishing

**Limitation**: This repository is named `xyops-integrations`, not `terraform-provider-xyops`. The Terraform Registry requires repo names matching `terraform-provider-{NAME}` for auto-detection.

This project will work via `dev_overrides`, but will not auto-publish to the public registry without:
1. Renaming the repo to `terraform-provider-xyops`, or
2. Extracting the provider into its own repo

## License

MIT License — see LICENSE file.
