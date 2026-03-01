# statuspage_component

Manages a Statuspage component.

## Example Usage

```hcl
resource "statuspage_component" "api" {
  page_id     = "abc123def456"
  name        = "API"
  description = "API service"
  status      = "operational"
  showcase    = true
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `name` (Required) - Display name for the component.
- `description` (Optional) - More detailed description of the component.
- `group_id` (Optional) - Component group identifier.
- `position` (Optional) - Order the component will appear on the page.
- `status` (Optional) - Status of the component. One of `operational`, `under_maintenance`, `degraded_performance`, `partial_outage`, `major_outage`.
- `showcase` (Optional) - Should this component be showcased.
- `only_show_if_degraded` (Optional) - Requires a special feature flag to be enabled.
- `start_date` (Optional) - The date this component started being used (format: `YYYY-MM-DD`).

## Attribute Reference

- `id` - Component identifier.
- `automation_email` - Automation email address.

## Import

```shell
terraform import statuspage_component.example page_id/component_id
```
