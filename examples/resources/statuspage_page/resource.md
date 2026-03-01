# statuspage_page

Manages a Statuspage page.

> **Note:** Pages cannot be created or deleted via the API. Use `terraform import` to manage existing pages. The `id` attribute must be set in the configuration.

## Example Usage

```hcl
resource "statuspage_page" "main" {
  id                       = "abc123def456"
  name                     = "My Company Status"
  page_description         = "Current status of My Company services"
  subdomain                = "status-mycompany"
  time_zone                = "Asia/Tokyo"
  allow_page_subscribers   = true
  allow_email_subscribers  = true
  allow_sms_subscribers    = false
  allow_rss_atom_feeds     = true
}
```

## Argument Reference

- `id` (Required) - Page identifier. Changing this forces a new resource.
- `name` (Optional) - Page name.
- `page_description` (Optional) - Page description.
- `headline` (Optional) - Page headline.
- `branding` (Optional) - Branding level.
- `subdomain` (Optional) - Page subdomain.
- `domain` (Optional) - Custom domain.
- `url` (Optional) - Page URL.
- `support_url` (Optional) - Support URL.
- `hidden_from_search` (Optional) - Hide page from search engines.
- `allow_page_subscribers` (Optional) - Allow page-level subscribers.
- `allow_incident_subscribers` (Optional) - Allow incident-level subscribers.
- `allow_email_subscribers` (Optional) - Allow email subscribers.
- `allow_sms_subscribers` (Optional) - Allow SMS subscribers.
- `allow_rss_atom_feeds` (Optional) - Allow RSS/Atom feeds.
- `allow_webhook_subscribers` (Optional) - Allow webhook subscribers.
- `notifications_from_email` (Optional) - Sender email for notifications.
- `notifications_email_footer` (Optional) - Footer text for notification emails.
- `time_zone` (Optional) - Timezone for the page.
- `city` (Optional) - City.
- `state` (Optional) - State.
- `country` (Optional) - Country.

## Import

```shell
terraform import statuspage_page.example page_id
```
