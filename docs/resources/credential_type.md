---
page_title: "aap_credential_type Resource - AAP Provider"
subcategory: ""
description: |-
  Manages a custom credential type in Ansible Automation Platform.
---

# aap_credential_type (Resource)

Manages a custom credential type in Ansible Automation Platform 2.5.

Custom credential types allow you to define new credential schemas for integration with external systems.

## Example Usage

```terraform
resource "aap_credential_type" "api_token" {
  name        = "API Token"
  description = "Custom credential for API authentication"
  kind        = "cloud"
  
  inputs = jsonencode({
    fields = [
      {
        id    = "api_token"
        type  = "string"
        label = "API Token"
        secret = true
      },
      {
        id    = "api_url"
        type  = "string"
        label = "API URL"
      }
    ]
    required = ["api_token", "api_url"]
  })
  
  injectors = jsonencode({
    env = {
      API_TOKEN = "{{ api_token }}"
      API_URL   = "{{ api_url }}"
    }
  })
}
```

## Argument Reference

### Required

- `name` (String) - Name of the credential type.
- `kind` (String) - Kind of credential: `"cloud"` or `"net"`.

### Optional

- `description` (String) - Description of the credential type.
- `inputs` (String) - Input field schema in JSON format.
- `injectors` (String) - Environment variable/file injector configuration in JSON format.

## Attribute Reference

- `id` - The ID of the credential type.

## Import

```shell
terraform import aap_credential_type.example 1
```
