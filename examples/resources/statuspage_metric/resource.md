# statuspage_metric

Manages a Statuspage metric.

## Example Usage

```hcl
resource "statuspage_metric" "response_time" {
  page_id             = "abc123def456"
  metrics_provider_id = statuspage_metrics_provider.datadog.id
  name                = "API Response Time"
  metric_identifier   = "api.response_time"
  display             = true
  tooltip_description = "Average API response time in milliseconds"
  suffix              = "ms"
  decimal_places      = 2
  y_axis_min          = 0
  y_axis_max          = 1000
  y_axis_hidden       = false
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `metrics_provider_id` (Required) - Metrics provider identifier. Changing this forces a new resource.
- `name` (Required) - Metric display name.
- `metric_identifier` (Optional) - Identifier of the metric in the metrics provider.
- `display` (Optional) - Whether to display this metric on the status page.
- `tooltip_description` (Optional) - Tooltip description shown on hover.
- `y_axis_min` (Optional) - Minimum value for the Y-axis.
- `y_axis_max` (Optional) - Maximum value for the Y-axis.
- `y_axis_hidden` (Optional) - Whether to hide the Y-axis.
- `suffix` (Optional) - Suffix to display after the metric value.
- `decimal_places` (Optional) - Number of decimal places to display.

## Attribute Reference

- `id` - Metric identifier.

## Import

```shell
terraform import statuspage_metric.example page_id/metric_id
```
