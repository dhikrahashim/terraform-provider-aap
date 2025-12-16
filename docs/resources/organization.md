---
page_title: "aap_organization Resource - AAP Provider"
subcategory: ""
description: |-
  Manages an organization in Ansible Automation Platform.
---

# aap_organization (Resource)

Manages an organization in Ansible Automation Platform 2.5.

Organizations are the highest level grouping of objects in AAP, including inventories, projects, job templates, and users.

## Example Usage

```terraform
resource "aap_organization" "example" {
  name        = "Production"
  description = "Production environment organization"
  max_hosts   = 100
}
```

## Argument Reference

The following arguments are supported:

### Required

- `name` (String) - The name of the organization. Must be unique.

### Optional

- `description` (String) - Description of the organization.
- `max_hosts` (Number) - Maximum number of hosts allowed to be managed by this organization. `0` means unlimited.
- `custom_virtualenv` (String) - Local absolute file path containing a custom Python virtualenv to use.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the organization.

## Import

Organizations can be imported using their ID:

```shell
terraform import aap_organization.example 1
```
