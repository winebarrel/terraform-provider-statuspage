# statuspage_page_access_user

Manages a Statuspage page access user.

## Example Usage

```hcl
resource "statuspage_page_access_user" "example" {
  page_id        = "abc123def456"
  external_login = "user@example.com"
  external_email = "user@example.com"
  component_ids  = [statuspage_component.api.id]
  metric_ids     = [statuspage_metric.response_time.id]
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `external_login` (Required) - External login or username.
- `external_email` (Optional) - External email address.
- `component_ids` (Optional) - List of component IDs the user can access.
- `metric_ids` (Optional) - List of metric IDs the user can access.

## Attribute Reference

- `id` - Page access user identifier.

## Import

```shell
terraform import statuspage_page_access_user.example page_id/page_access_user_id
```
