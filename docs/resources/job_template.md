---
page_title: "aap_job_template Resource - AAP Provider"
subcategory: ""
description: |-
  Manages a job template in Ansible Automation Platform.
---

# aap_job_template (Resource)

Manages a job template in Ansible Automation Platform 2.5.

A job template is a definition and set of parameters for running an Ansible job.

## Example Usage

### Basic Job Template

```terraform
resource "aap_job_template" "deploy" {
  name         = "Deploy Application"
  job_type     = "run"
  inventory_id = aap_inventory.production.id
  project_id   = "1"
  playbook     = "deploy.yml"
}
```

### Job Template with Extra Variables

```terraform
resource "aap_job_template" "configure" {
  name         = "Configure Servers"
  job_type     = "run"
  inventory_id = aap_inventory.production.id
  project_id   = "1"
  playbook     = "configure.yml"
  verbosity    = 2
  forks        = 10
  limit        = "webservers"
  extra_vars   = jsonencode({
    environment = "production"
    version     = "1.2.3"
  })
}
```

### Check Mode Job Template

```terraform
resource "aap_job_template" "check" {
  name         = "Check Configuration"
  job_type     = "check"  # Dry run mode
  inventory_id = aap_inventory.production.id
  project_id   = "1"
  playbook     = "configure.yml"
}
```

## Argument Reference

### Required

- `name` (String) - Name of the job template.
- `job_type` (String) - Type of job. Valid values: `"run"`, `"check"`.
- `inventory_id` (String) - ID of the inventory to use.
- `project_id` (String) - ID of the project containing the playbook.
- `playbook` (String) - Name of the playbook to run.

### Optional

- `description` (String) - Description of the job template.
- `forks` (Number) - Number of parallel processes to use. Default: `0` (use Ansible default).
- `limit` (String) - Host pattern to limit execution.
- `verbosity` (Number) - Verbosity level (0-5). Default: `0`.
- `extra_vars` (String) - Extra variables in JSON or YAML format.

## Attribute Reference

- `id` - The ID of the job template.

## Import

```shell
terraform import aap_job_template.example 1
```
