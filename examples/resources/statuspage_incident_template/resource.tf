import {
  to = statuspage_incident_template.outage
  id = "abc123def456/template-id-1"
}

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
