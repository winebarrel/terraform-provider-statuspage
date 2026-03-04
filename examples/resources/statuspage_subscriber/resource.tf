resource "statuspage_subscriber" "email" {
  page_id       = "abc123def456"
  email         = "user@example.com"
  component_ids = [statuspage_component.api.id]

  skip_confirmation_notification = false
}
