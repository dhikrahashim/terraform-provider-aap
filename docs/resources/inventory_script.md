---
page_title: "aap_inventory_script Resource - AAP Provider"
subcategory: ""
description: |-
  Manages a custom inventory script in Ansible Automation Platform.
---

# aap_inventory_script (Resource)

Manages a custom inventory script in Ansible Automation Platform 2.5.

Inventory scripts are custom scripts that generate dynamic inventory data.

## Example Usage

```terraform
resource "aap_inventory_script" "example" {
  name            = "Custom AWS Inventory"
  organization_id = aap_organization.example.id
  description     = "Custom script to fetch AWS inventory"
  
  script = <<-EOF
    #!/usr/bin/env python
    import json
    
    inventory = {
        "all": {
            "hosts": ["host1.example.com", "host2.example.com"]
        }
    }
    
    print(json.dumps(inventory))
  EOF
}
```

## Argument Reference

### Required

- `name` (String) - Name of the inventory script.
- `organization_id` (String) - Organization ID.
- `script` (String) - The inventory script content (Python or shell script).

### Optional

- `description` (String) - Description of the inventory script.

## Attribute Reference

- `id` - The ID of the inventory script.

## Import

```shell
terraform import aap_inventory_script.example 1
```
