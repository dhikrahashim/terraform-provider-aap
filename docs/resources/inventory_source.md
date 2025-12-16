---
page_title: "aap_inventory_source Resource - AAP Provider"
subcategory: ""
description: |-
  Manages an inventory source in Ansible Automation Platform.
---

# aap_inventory_source (Resource)

Manages an inventory source in Ansible Automation Platform 2.5.

Inventory sources dynamically populate inventories from external systems like cloud providers, virtualization platforms, or custom scripts.

## Example Usage

### SCM Source (from Project)

```terraform
resource "aap_inventory_source" "from_project" {
  name              = "Hosts from Git"
  inventory_id      = aap_inventory.example.id
  source            = "scm"
  source_project_id = aap_project.example.id
  source_path       = "inventory/hosts.yml"
  update_on_launch  = true
}
```

### AWS EC2 Source

```terraform
resource "aap_inventory_source" "aws" {
  name           = "AWS Instances"
  inventory_id   = aap_inventory.example.id
  source         = "ec2"
  credential_id  = aap_credential_cloud.aws.id
  source_vars    = jsonencode({
    regions = ["us-east-1", "us-west-2"]
  })
}
```

## Argument Reference

### Required

- `name` (String) - Name of the inventory source.
- `inventory_id` (String) - Inventory ID.
- `source` (String) - Source type: `"scm"`, `"ec2"`, `"gce"`, `"azure_rm"`, `"vmware"`, `"satellite6"`, `"openstack"`, `"rhv"`, `"controller"`, `"file"`.

### Optional

- `description` (String) - Description.
- `source_path` (String) - Path to inventory file within a project (for `scm` source).
- `source_vars` (String) - Source-specific variables in YAML/JSON format.
- `credential_id` (String) - Cloud/network credential ID.
- `source_project_id` (String) - Project ID containing inventory file (for `scm` source).
- `update_on_launch` (Boolean) - Update inventory when a job is launched.
- `update_cache_timeout` (Number) - Cache timeout for updates.
- `overwrite` (Boolean) - Overwrite local groups and hosts.
- `overwrite_vars` (Boolean) - Overwrite local variables.

## Attribute Reference

- `id` - The ID of the inventory source.

## Import

```shell
terraform import aap_inventory_source.example 1
```
