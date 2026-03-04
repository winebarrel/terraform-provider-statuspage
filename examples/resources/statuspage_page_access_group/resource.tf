resource "statuspage_page_access_group" "engineering" {
  page_id              = "abc123def456"
  name                 = "Engineering Team"
  external_identifier  = "eng-team"
  component_ids        = [statuspage_component.api.id]
  metric_ids           = [statuspage_metric.response_time.id]
  page_access_user_ids = [statuspage_page_access_user.example.id]
}
