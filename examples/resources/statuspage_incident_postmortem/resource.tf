resource "statuspage_incident_postmortem" "example" {
  page_id     = "abc123def456"
  incident_id = statuspage_incident.outage.id

  body               = "## Root Cause\nThe outage was caused by a database connection pool exhaustion.\n\n## Resolution\nWe increased the connection pool size and added monitoring alerts."
  notify_subscribers = true
  notify_twitter     = false
}
