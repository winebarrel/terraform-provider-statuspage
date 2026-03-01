# statuspage_page_access_group

Manages a Statuspage page access group.

## Example Usage

```hcl
resource "statuspage_page_access_group" "engineering" {
  page_id             = "abc123def456"
  name                = "Engineering Team"
  external_identifier = "eng-team"
  component_ids       = [statuspage_component.api.id]
  metric_ids          = [statuspage_metric.response_time.id]
  page_access_user_ids = [statuspage_page_access_user.example.id]
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `name` (Required) - Group name.
- `external_identifier` (Optional) - External identifier for the group.
- `component_ids` (Optional) - List of component IDs the group can access.
- `metric_ids` (Optional) - List of metric IDs the group can access.
- `page_access_user_ids` (Optional) - List of page access user IDs in this group.

## Attribute Reference

- `id` - Page access group identifier.

## Import

```shell
terraform import statuspage_page_access_group.example page_id/page_access_group_id
```
