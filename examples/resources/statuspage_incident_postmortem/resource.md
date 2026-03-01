# statuspage_incident_postmortem

Manages a Statuspage incident postmortem.

## Example Usage

```hcl
resource "statuspage_incident_postmortem" "example" {
  page_id     = "abc123def456"
  incident_id = statuspage_incident.outage.id

  body               = "## Root Cause\nThe outage was caused by a database connection pool exhaustion.\n\n## Resolution\nWe increased the connection pool size and added monitoring alerts."
  notify_subscribers = true
  notify_twitter     = false
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `incident_id` (Required) - Incident identifier. Changing this forces a new resource.
- `body` (Optional) - Postmortem body (supports markdown).
- `body_draft` (Optional) - Postmortem draft body.
- `notify_subscribers` (Optional) - Notify subscribers when postmortem is published.
- `notify_twitter` (Optional) - Notify on Twitter when postmortem is published.

## Import

```shell
terraform import statuspage_incident_postmortem.example page_id/incident_id
```
