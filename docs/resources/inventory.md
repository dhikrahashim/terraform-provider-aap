---
page_title: "aap_inventory Resource - AAP Provider"
subcategory: ""
description: |-
  Manages an inventory in Ansible Automation Platform.
---

# aap_inventory (Resource)

Manages an inventory in Ansible Automation Platform 2.5.

An inventory is a collection of hosts against which jobs may be launched.

## Example Usage

### Standard Inventory

```terraform
resource "aap_inventory" "example" {
  name            = "Production Servers"
  organization_id = aap_organization.example.id
  description     = "Inventory for production servers"
}
```

### Smart Inventory

```terraform
resource "aap_inventory" "smart" {
  name            = "Linux Servers"
  organization_id = aap_organization.example.id
  kind            = "smart"
  host_filter     = "ansible_os_family=RedHat"
}
```

### Inventory with Variables

```terraform
resource "aap_inventory" "with_vars" {
  name            = "Web Servers"
  organization_id = aap_organization.example.id
  variables       = jsonencode({
    ansible_user = "deploy"
    http_port    = 8080
  })
}
```

## Argument Reference

### Required

- `name` (String) - The name of the inventory.
- `organization_id` (String) - The ID of the organization containing this inventory.

### Optional

- `description` (String) - Description of the inventory.
- `kind` (String) - Kind of inventory. Empty string for standard inventory, `"smart"` for smart inventory.
- `host_filter` (String) - Filter for smart inventories. Only applicable when `kind = "smart"`.
- `variables` (String) - Inventory variables in JSON or YAML format.

## Attribute Reference

- `id` - The ID of the inventory.

## Import

```shell
terraform import aap_inventory.example 1
```
