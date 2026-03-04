import {
  to = statuspage_metric.response_time
  id = "abc123def456/metric-id-1"
}

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
