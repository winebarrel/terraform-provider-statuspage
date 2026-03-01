# statuspage_user_permissions

Manages Statuspage user permissions.

> **Note:** User permissions cannot be deleted via the API. Removing this resource from your configuration will only remove it from Terraform state.

## Example Usage

```hcl
resource "statuspage_user_permissions" "example" {
  organization_id = "org123"
  user_id         = statuspage_user.example.id

  pages = {
    "page_id_1" = "manage"
    "page_id_2" = "view"
  }
}
```

## Argument Reference

- `organization_id` (Required) - Organization identifier. Changing this forces a new resource.
- `user_id` (Required) - User identifier. Changing this forces a new resource.
- `pages` (Required) - Map of page IDs to permission levels.

## Import

```shell
terraform import statuspage_user_permissions.example organization_id/user_id
```
