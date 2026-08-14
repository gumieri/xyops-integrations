# Terraform Provider Example

## Local Development Setup

1. Add to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "gumieri/xyops" = "/path/to/xyops-integrations"
  }
  direct {}
}
```

2. Run `terraform init` in this directory

## Example Usage

```hcl
terraform {
  required_providers {
    xyops = {
      source = "gumieri/xyops"
    }
  }
}

provider "xyops" {
  base_url = var.xyops_base_url
  api_key  = var.xyops_api_key
}

variable "xyops_base_url" {
  type    = string
  default = "https://xyops.example.com"
}

variable "xyops_api_key" {
  type      = string
  sensitive = true
}

# Tag
resource "xyops_tag" "backup" {
  title = "backup"
}

# Group
resource "xyops_group" "vps" {
  title          = "vps"
  hostname_match = "*-vps-*"
}

# Category
resource "xyops_category" "maintenance" {
  title   = "maintenance"
  enabled = true
}

# Channel
resource "xyops_channel" "oncall" {
  title = "oncall"
}

# Secret (field values are encrypted and write-only)
resource "xyops_secret" "db_credentials" {
  title   = "db-credentials"
  enabled = true
  fields = [
    { name = "host", value = "db.internal" },
    { name = "password", value = var.db_password },
  ]
}

# Bucket
resource "xyops_bucket" "artifacts" {
  title   = "artifacts"
  enabled = true
  notes   = "Build artifacts"
}

# Ticket
resource "xyops_ticket" "example" {
  subject = "Example ticket"
  body    = "Created by Terraform"
  status  = "open"
  tags    = ["automated"]
}

# Data sources
data "xyops_channels" "all" {}
data "xyops_secrets" "all" {}
data "xyops_buckets" "all" {}
data "xyops_tickets" "all" {}

```