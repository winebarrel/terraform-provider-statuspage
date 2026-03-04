import {
  to = statuspage_incident.outage
  id = "abc123def456/incident-id-1"
}

resource "statuspage_incident" "outage" {
  page_id         = "abc123def456"
  name            = "API Degraded Performance"
  status          = "investigating"
  body            = "We are currently investigating degraded API performance."
  impact_override = "minor"
}
