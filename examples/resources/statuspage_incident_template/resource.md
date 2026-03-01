# statuspage_incident_template

Manages a Statuspage incident template.

## Example Usage

```hcl
resource "statuspage_incident_template" "outage" {
  page_id       = "abc123def456"
  name          = "Service Outage Template"
  title         = "Service Outage"
  body          = "We are currently experiencing a service outage. Our team is investigating."
  update_status = "investigating"
  component_ids = [statuspage_component.api.id]

  should_tweet              = false
  should_send_notifications = true
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `name` (Required) - Template name.
- `title` (Optional) - Incident title to use when creating an incident from this template.
- `body` (Optional) - Incident body to use when creating an incident from this template.
- `group_id` (Optional) - Group identifier.
- `update_status` (Optional) - Status to set on the incident. One of `investigating`, `identified`, `monitoring`, `resolved`, `scheduled`, `in_progress`, `verifying`, `completed`.
- `should_tweet` (Optional) - Tweet when incident is created from this template.
- `should_send_notifications` (Optional) - Send notifications when incident is created from this template.
- `component_ids` (Optional) - List of component IDs to associate with incidents created from this template.

## Attribute Reference

- `id` - Incident template identifier.

## Import

```shell
terraform import statuspage_incident_template.example page_id/template_id
```
