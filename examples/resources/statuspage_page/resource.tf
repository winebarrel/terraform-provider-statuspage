resource "statuspage_page" "main" {
  id                      = "abc123def456"
  name                    = "My Company Status"
  page_description        = "Current status of My Company services"
  subdomain               = "status-mycompany"
  time_zone               = "Asia/Tokyo"
  allow_page_subscribers  = true
  allow_email_subscribers = true
  allow_sms_subscribers   = false
  allow_rss_atom_feeds    = true
}
