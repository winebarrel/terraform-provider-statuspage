resource "statuspage_component" "api" {
  page_id     = "abc123def456"
  name        = "API"
  description = "API service"
  status      = "operational"
  showcase    = true
}
