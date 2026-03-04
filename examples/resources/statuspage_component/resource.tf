import {
  to = statuspage_component.api
  id = "abc123def456/component-id-1"
}

resource "statuspage_component" "api" {
  page_id     = "abc123def456"
  name        = "API"
  description = "API service"
  status      = "operational"
  showcase    = true
}
