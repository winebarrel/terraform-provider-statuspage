resource "statuspage_component" "api" {
  page_id = "abc123def456"
  name    = "API"
}

resource "statuspage_component" "web" {
  page_id = "abc123def456"
  name    = "Web App"
}

import {
  to = statuspage_component_group.services
  id = "abc123def456/component-group-id-1"
}

resource "statuspage_component_group" "services" {
  page_id     = "abc123def456"
  name        = "Core Services"
  description = "Core platform services"
  components = [
    statuspage_component.api.id,
    statuspage_component.web.id,
  ]
}
