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
```