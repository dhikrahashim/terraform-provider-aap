---
page_title: "aap_credential_scm Resource - AAP Provider"
subcategory: ""
description: |-
  Manages an SCM credential in Ansible Automation Platform.
---

# aap_credential_scm (Resource)

Manages a Source Control (SCM) credential in Ansible Automation Platform 2.5.

SCM credentials are used to access private Git, Mercurial, or Subversion repositories.

## Example Usage

### Username/Password

```terraform
resource "aap_credential_scm" "github" {
  name            = "GitHub Credentials"
  organization_id = aap_organization.example.id
  username        = "git-user"
  password        = "personal-access-token"
}
```

### SSH Key

```terraform
resource "aap_credential_scm" "git_ssh" {
  name            = "Git SSH Key"
  organization_id = aap_organization.example.id
  username        = "git"
  ssh_key_data    = file("~/.ssh/id_rsa")
}
```

## Argument Reference

### Required

- `name` (String) - Name of the credential.
- `organization_id` (String) - Organization ID.

### Optional

- `description` (String) - Description of the credential.
- `username` (String) - SCM username.
- `password` (String, Sensitive) - SCM password or personal access token.
- `ssh_key_data` (String, Sensitive) - Private SSH key.
- `ssh_key_unlock` (String, Sensitive) - Passphrase for encrypted SSH key.

## Attribute Reference

- `id` - The ID of the credential.

## Import

```shell
terraform import aap_credential_scm.example 1
```
