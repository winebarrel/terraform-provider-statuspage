# statuspage_user

Manages a Statuspage organization user.

> **Note:** This resource is not available for organizations using Atlassian accounts. Updates are not supported; changing any attribute will force recreation of the resource.

## Example Usage

```hcl
resource "statuspage_user" "example" {
  organization_id = "org123"
  email           = "newuser@example.com"
  password        = "securepassword123"
  first_name      = "John"
  last_name       = "Doe"
}
```

## Argument Reference

- `organization_id` (Required) - Organization identifier. Changing this forces a new resource.
- `email` (Required) - User email address. Changing this forces a new resource.
- `password` (Required, Sensitive) - User password. Changing this forces a new resource.
- `first_name` (Optional) - User first name.
- `last_name` (Optional) - User last name.

## Attribute Reference

- `id` - User identifier.

## Import

```shell
terraform import statuspage_user.example organization_id/user_id
```
