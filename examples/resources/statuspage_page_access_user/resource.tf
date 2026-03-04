resource "statuspage_page_access_user" "example" {
  page_id        = "abc123def456"
  external_login = "user@example.com"
  external_email = "user@example.com"
  component_ids  = [statuspage_component.api.id]
  metric_ids     = [statuspage_metric.response_time.id]
}
