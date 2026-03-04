resource "statuspage_status_embed_config" "example" {
  page_id                      = "abc123def456"
  position                     = "bottom-right"
  incident_background_color    = "#FF6347"
  incident_text_color          = "#FFFFFF"
  maintenance_background_color = "#FFA500"
  maintenance_text_color       = "#FFFFFF"
}
