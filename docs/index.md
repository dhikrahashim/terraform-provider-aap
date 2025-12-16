---
page_title: "AAP Provider"
subcategory: ""
description: |-
  Terraform provider for managing Ansible Automation Platform (AAP) 2.5 resources.
---

# AAP Provider

The AAP provider allows you to manage resources in Red Hat Ansible Automation Platform 2.5 via Terraform.

This provider uses the new AAP 2.5 API path (`/api/controller/v2/`) through the Platform Gateway.

## Example Usage

```terraform
terraform {
  required_providers {
    aap = {
      source  = "dhikrahashim/aap"
      version = "~> 0.1.0"
    }
  }
}

provider "aap" {
  host     = "https://aap.example.com"
  username = "admin"
  password = "password"
  insecure = true
}
```

## Authentication

The provider supports two authentication methods:

### Basic Authentication
```terraform
provider "aap" {
  host     = "https://aap.example.com"
  username = "admin"
  password = "password"
}
```

### OAuth2 Token
```terraform
provider "aap" {
  host  = "https://aap.example.com"
  token = "your_oauth_token"
}
```

## Configuration via Environment Variables

You can configure the provider using environment variables:

| Variable | Description |
|----------|-------------|
| `AAP_HOST` | AAP Controller URL |
| `AAP_USERNAME` | Username for authentication |
| `AAP_PASSWORD` | Password for authentication |
| `AAP_TOKEN` | OAuth2 token (alternative to username/password) |

## Schema

### Optional

- `host` (String) - The URI of the AAP Controller (e.g., `https://aap.example.com`)
- `username` (String) - Username for authentication
- `password` (String, Sensitive) - Password for authentication
- `token` (String, Sensitive) - OAuth2 token for authentication
- `insecure` (Boolean) - Skip TLS certificate verification (default: `false`)
