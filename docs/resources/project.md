---
page_title: "aap_project Resource - AAP Provider"
subcategory: ""
description: |-
  Manages a project in Ansible Automation Platform.
---

# aap_project (Resource)

Manages a project in Ansible Automation Platform 2.5.

Projects are logical collections of Ansible playbooks, sourced from version control systems.

## Example Usage

### Git Project

```terraform
resource "aap_project" "example" {
  name              = "My Playbooks"
  organization_id   = aap_organization.example.id
  scm_type          = "git"
  scm_url           = "https://github.com/example/playbooks.git"
  scm_branch        = "main"
  scm_update_on_launch = true
}
```

### Project with SCM Credential

```terraform
resource "aap_project" "private_repo" {
  name              = "Private Playbooks"
  organization_id   = aap_organization.example.id
  scm_type          = "git"
  scm_url           = "git@github.com:example/private-playbooks.git"
  scm_branch        = "main"
  scm_credential_id = aap_credential_scm.git_ssh.id
}
```

## Argument Reference

### Required

- `name` (String) - Name of the project.
- `organization_id` (String) - Organization ID.
- `scm_type` (String) - SCM type: `""` (manual), `"git"`, `"hg"`, `"svn"`.

### Optional

- `description` (String) - Description of the project.
- `scm_url` (String) - SCM repository URL.
- `scm_branch` (String) - Branch, tag, or commit to checkout.
- `scm_credential_id` (String) - SCM credential ID.
- `scm_clean` (Boolean) - Clean the repository before syncing.
- `scm_delete_on_update` (Boolean) - Delete local modifications before updating.
- `scm_update_on_launch` (Boolean) - Update project when a job is launched.
- `scm_update_cache_timeout` (Number) - Cache timeout for SCM updates.
- `local_path` (String) - Local path for manual projects.

## Attribute Reference

- `id` - The ID of the project.

## Import

```shell
terraform import aap_project.example 1
```
