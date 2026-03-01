# statuspage_status_embed_config

Manages a Statuspage status embed configuration.

> **Note:** Status embed configuration cannot be created or deleted via the API. Use `terraform import` to manage existing configurations.

## Example Usage

```hcl
resource "statuspage_status_embed_config" "example" {
  page_id                      = "abc123def456"
  position                     = "bottom-right"
  incident_background_color    = "#FF6347"
  incident_text_color          = "#FFFFFF"
  maintenance_background_color = "#FFA500"
  maintenance_text_color       = "#FFFFFF"
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `position` (Optional) - Position of the status embed widget.
- `incident_background_color` (Optional) - Background color for incident status.
- `incident_text_color` (Optional) - Text color for incident status.
- `maintenance_background_color` (Optional) - Background color for maintenance status.
- `maintenance_text_color` (Optional) - Text color for maintenance status.

## Import

```shell
terraform import statuspage_status_embed_config.example page_id
```
