# statuspage_metrics_provider

Manages a Statuspage metrics provider.

## Example Usage

### Datadog

```hcl
resource "statuspage_metrics_provider" "datadog" {
  page_id         = "abc123def456"
  type            = "Datadog"
  api_key         = "your-datadog-api-key"
  application_key = "your-datadog-app-key"
}
```

### Self (Custom Metrics)

```hcl
resource "statuspage_metrics_provider" "custom" {
  page_id = "abc123def456"
  type    = "Self"
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `type` (Required) - Provider type. One of `Pingdom`, `NewRelic`, `Datadog`, `Self`, `Librato`, `CalabashCustom`.
- `email` (Optional, Sensitive) - Email for the metrics provider.
- `metric_base_uri` (Optional) - Metric base URI.
- `api_key` (Optional, Sensitive) - API key for the metrics provider. Not returned by the API after creation.
- `api_token` (Optional, Sensitive) - API token for the metrics provider. Not returned by the API after creation.
- `application_key` (Optional, Sensitive) - Application key for the metrics provider. Not returned by the API after creation.

## Attribute Reference

- `id` - Metrics provider identifier.

## Import

```shell
terraform import statuspage_metrics_provider.example page_id/metrics_provider_id
```
