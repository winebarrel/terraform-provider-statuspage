resource "statuspage_metrics_provider" "datadog" {
  page_id         = "abc123def456"
  type            = "Datadog"
  api_key         = "your-datadog-api-key"
  application_key = "your-datadog-app-key"
}
