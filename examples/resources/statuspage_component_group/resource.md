# statuspage_component_group

Manages a Statuspage component group.

## Example Usage

```hcl
resource "statuspage_component" "api" {
  page_id = "abc123def456"
  name    = "API"
}

resource "statuspage_component" "web" {
  page_id = "abc123def456"
  name    = "Web App"
}

resource "statuspage_component_group" "services" {
  page_id     = "abc123def456"
  name        = "Core Services"
  description = "Core platform services"
  components  = [
    statuspage_component.api.id,
    statuspage_component.web.id,
  ]
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `name` (Required) - Display name for the component group.
- `components` (Required) - List of component IDs that belong to this group.
- `description` (Optional) - Description of the component group.
- `position` (Optional) - Order the component group will appear on the page.

## Attribute Reference

- `id` - Component group identifier.

## Import

```shell
terraform import statuspage_component_group.example page_id/component_group_id
```
