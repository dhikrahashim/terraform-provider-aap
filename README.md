# Terraform Provider for Ansible Automation Platform (AAP) 2.5

[![Release](https://github.com/dhikrahashim/terraform-provider-aap/actions/workflows/release.yml/badge.svg)](https://github.com/dhikrahashim/terraform-provider-aap/actions/workflows/release.yml)

This provider allows you to manage resources in Red Hat Ansible Automation Platform (AAP) 2.5 via Terraform.

**Key Feature**: Uses the new AAP 2.5 API path (`/api/controller/v2/`) via the Platform Gateway.

## Requirements

- Terraform >= 1.0
- Go >= 1.22 (to build the provider plugin)
- AAP 2.5 instance

## Resources Supported

| Resource | Description |
|----------|-------------|
| `aap_organization` | Manage organizations |
| `aap_inventory` | Manage inventories (standard and smart) |
| `aap_job_template` | Manage job templates |

## Building The Provider

```bash
git clone https://github.com/dhikrahashim/terraform-provider-aap.git
cd terraform-provider-aap
go mod tidy
go build -o terraform-provider-aap
```

## Installing Locally

Add to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "dhikrahashim/aap" = "/path/to/directory/containing/binary"
  }
  direct {}
}
```

## Usage Example

```hcl
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
  # Or use OAuth token:
  # token = "your_oauth_token"
  insecure = true  # Skip TLS verification (dev only)
}

# Create an organization
resource "aap_organization" "example" {
  name        = "Terraform Managed Org"
  description = "Created via Terraform"
  max_hosts   = 100
}

# Create an inventory
resource "aap_inventory" "example" {
  name            = "Production Servers"
  organization_id = aap_organization.example.id
  description     = "Production environment inventory"
}

# Create a job template
resource "aap_job_template" "example" {
  name         = "Deploy Application"
  job_type     = "run"
  inventory_id = aap_inventory.example.id
  project_id   = "1"  # Your project ID
  playbook     = "deploy.yml"
  verbosity    = 1
}
```

## Environment Variables

You can also configure the provider using environment variables:

| Variable | Description |
|----------|-------------|
| `AAP_HOST` | AAP Controller URL |
| `AAP_USERNAME` | Username for authentication |
| `AAP_PASSWORD` | Password for authentication |
| `AAP_TOKEN` | OAuth2 token (alternative to username/password) |

## Publishing to Terraform Registry

See [PUBLISHING.md](PUBLISHING.md) for step-by-step instructions.

## License

MIT
