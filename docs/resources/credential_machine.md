---
page_title: "aap_credential_machine Resource - AAP Provider"
subcategory: ""
description: |-
  Manages a machine credential in Ansible Automation Platform.
---

# aap_credential_machine (Resource)

Manages a machine (SSH) credential in Ansible Automation Platform 2.5.

Machine credentials are used to authenticate to managed hosts via SSH.

## Example Usage

### Password Authentication

```terraform
resource "aap_credential_machine" "example" {
  name            = "Linux Servers"
  organization_id = aap_organization.example.id
  username        = "ansible"
  password        = "secret123"
}
```

### SSH Key Authentication

```terraform
resource "aap_credential_machine" "ssh_key" {
  name            = "SSH Key Auth"
  organization_id = aap_organization.example.id
  username        = "ansible"
  ssh_key_data    = file("~/.ssh/id_rsa")
}
```

### With Privilege Escalation

```terraform
resource "aap_credential_machine" "with_sudo" {
  name            = "Sudo Access"
  organization_id = aap_organization.example.id
  username        = "ansible"
  ssh_key_data    = file("~/.ssh/id_rsa")
  become_method   = "sudo"
  become_username = "root"
  become_password = "rootpassword"
}
```

## Argument Reference

### Required

- `name` (String) - Name of the credential.
- `organization_id` (String) - Organization ID.

### Optional

- `description` (String) - Description of the credential.
- `username` (String) - SSH username.
- `password` (String, Sensitive) - SSH password.
- `ssh_key_data` (String, Sensitive) - Private SSH key.
- `ssh_public_key_data` (String) - Public SSH key.
- `ssh_key_unlock` (String, Sensitive) - Passphrase for encrypted SSH key.
- `become_method` (String) - Privilege escalation method: `sudo`, `su`, `pbrun`, `pfexec`, `dzdo`, `pmrun`, `runas`.
- `become_username` (String) - Privilege escalation username.
- `become_password` (String, Sensitive) - Privilege escalation password.

## Attribute Reference

- `id` - The ID of the credential.

## Import

```shell
terraform import aap_credential_machine.example 1
```
