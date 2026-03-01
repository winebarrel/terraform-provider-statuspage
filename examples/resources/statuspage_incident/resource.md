# statuspage_incident

Manages a Statuspage incident.

## Example Usage

### Basic Incident

```hcl
resource "statuspage_incident" "outage" {
  page_id = "abc123def456"
  name    = "API Degraded Performance"
  status  = "investigating"
  body    = "We are currently investigating degraded API performance."

  impact_override = "minor"
}
```

### Incident with Affected Components

```hcl
resource "statuspage_incident" "outage" {
  page_id        = "abc123def456"
  name           = "API Outage"
  status         = "identified"
  body           = "We have identified the root cause."
  impact_override = "major"
  component_ids  = [statuspage_component.api.id]

  components = {
    (statuspage_component.api.id) = "major_outage"
  }
}
```

### Scheduled Maintenance

```hcl
resource "statuspage_incident" "maintenance" {
  page_id        = "abc123def456"
  name           = "Database Maintenance"
  status         = "scheduled"
  body           = "Scheduled database maintenance window."
  scheduled_for  = "2025-01-15T02:00:00Z"
  scheduled_until = "2025-01-15T04:00:00Z"

  scheduled_remind_prior      = true
  scheduled_auto_in_progress  = true
  scheduled_auto_completed    = true
  deliver_notifications       = true
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `name` (Required) - Incident name.
- `status` (Optional) - Incident status. One of `investigating`, `identified`, `monitoring`, `resolved`, `scheduled`, `in_progress`, `verifying`, `completed`.
- `impact_override` (Optional) - Impact override. One of `none`, `minor`, `major`, `critical`, `maintenance`.
- `body` (Optional) - Initial incident update message.
- `component_ids` (Optional) - List of component IDs affected by this incident.
- `components` (Optional) - Map of component IDs to their status.
- `scheduled_for` (Optional) - Scheduled maintenance start time (ISO 8601).
- `scheduled_until` (Optional) - Scheduled maintenance end time (ISO 8601).
- `scheduled_remind_prior` (Optional) - Send reminder before scheduled maintenance.
- `scheduled_auto_in_progress` (Optional) - Automatically transition to in_progress at start.
- `scheduled_auto_completed` (Optional) - Automatically transition to completed at end.
- `auto_transition_to_maintenance_state` (Optional) - Automatically transition to maintenance state.
- `auto_transition_to_operational_state` (Optional) - Automatically transition to operational state.
- `auto_transition_deliver_notifications_at_start` (Optional) - Deliver notifications at maintenance start.
- `auto_transition_deliver_notifications_at_end` (Optional) - Deliver notifications at maintenance end.
- `deliver_notifications` (Optional) - Deliver notifications to subscribers. Defaults to `true`.
- `auto_tweet_at_beginning` (Optional) - Tweet at beginning of incident.
- `auto_tweet_on_creation` (Optional) - Tweet on incident creation.
- `auto_tweet_on_completion` (Optional) - Tweet on incident completion.
- `auto_tweet_one_hour_before` (Optional) - Tweet one hour before scheduled maintenance.

## Attribute Reference

- `id` - Incident identifier.
- `shortlink` - Incident shortlink URL.

## Import

```shell
terraform import statuspage_incident.example page_id/incident_id
```
