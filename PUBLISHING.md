# Publishing terraform-provider-aap to Terraform Registry

This guide provides step-by-step instructions for publishing your provider.

## Prerequisites

- GitHub account (dhikrahashim)
- Terraform Registry account (namespace: dhikrahashim)
- GPG key for signing releases

---

## Step 1: Generate GPG Key (if you don't have one)

```bash
# Generate a new GPG key
gpg --full-generate-key
# Choose: RSA and RSA, 4096 bits, no expiration

# List your keys to get the key ID
gpg --list-secret-keys --keyid-format=long

# Export your public key (replace KEY_ID)
gpg --armor --export KEY_ID > public.asc

# Export your private key (for GitHub secrets)
gpg --armor --export-secret-keys KEY_ID > private.asc
```

---

## Step 2: Add GPG Key to Terraform Registry

1. Go to https://registry.terraform.io/settings/gpg-keys
2. Click "Add GPG Key"
3. Paste the contents of `public.asc`
4. Save

---

## Step 3: Set Up GitHub Repository

1. Create a new repository: `terraform-provider-aap`
   - Go to https://github.com/new
   - Name: `terraform-provider-aap`
   - Make it **Public** (required for registry)

2. Push your code:
```bash
cd /Users/hashimabdulla/terraform-provider-aap

# Initialize git (if not already)
git init

# Add all files
git add .

# Commit
git commit -m "Initial commit: terraform-provider-aap for AAP 2.5"

# Add remote
git remote add origin https://github.com/dhikrahashim/terraform-provider-aap.git

# Push
git push -u origin main
```

---

## Step 4: Add GitHub Secrets

Go to your repository → Settings → Secrets and variables → Actions → New repository secret

Add these secrets:

| Secret Name | Value |
|-------------|-------|
| `GPG_PRIVATE_KEY` | Contents of `private.asc` (your GPG private key) |
| `PASSPHRASE` | Your GPG key passphrase (if any) |

---

## Step 5: Create and Push a Release Tag

```bash
# Create a tag
git tag v0.1.0

# Push the tag (this triggers the release workflow)
git push origin v0.1.0
```

The GitHub Action will automatically:
- Build binaries for all platforms
- Sign the release with your GPG key
- Create a GitHub Release with all artifacts

---

## Step 6: Register Provider in Terraform Registry

1. Go to https://registry.terraform.io/publish/provider
2. Sign in with GitHub
3. Select repository: `dhikrahashim/terraform-provider-aap`
4. Follow the prompts to complete registration

---

## Step 7: Verify Publication

After registration, your provider will be available at:
```
https://registry.terraform.io/providers/dhikrahashim/aap/latest
```

Users can use it with:
```hcl
terraform {
  required_providers {
    aap = {
      source  = "dhikrahashim/aap"
      version = "~> 0.1.0"
    }
  }
}
```

---

## Troubleshooting

### "No releases found"
- Ensure your tag starts with `v` (e.g., `v0.1.0`)
- Check GitHub Actions for errors

### "Invalid GPG signature"
- Ensure the GPG key in Registry matches the one used for signing
- Check that `GPG_PRIVATE_KEY` secret is correct

### Build fails
- Run `go build` locally first to verify
- Check Go version matches `go.mod`
