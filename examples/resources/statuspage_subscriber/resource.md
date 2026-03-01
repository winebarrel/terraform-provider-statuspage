# statuspage_subscriber

Manages a Statuspage subscriber.

## Example Usage

### Email Subscriber

```hcl
resource "statuspage_subscriber" "email" {
  page_id       = "abc123def456"
  email         = "user@example.com"
  component_ids = [statuspage_component.api.id]

  skip_confirmation_notification = false
}
```

### Webhook Subscriber

```hcl
resource "statuspage_subscriber" "webhook" {
  page_id  = "abc123def456"
  endpoint = "https://example.com/webhook"
}
```

### SMS Subscriber

```hcl
resource "statuspage_subscriber" "sms" {
  page_id       = "abc123def456"
  phone_number  = "+1234567890"
  phone_country = "US"
}
```

## Argument Reference

- `page_id` (Required) - Page identifier. Changing this forces a new resource.
- `email` (Optional) - Subscriber email address.
- `phone_number` (Optional) - Subscriber phone number.
- `phone_country` (Optional) - Phone country code.
- `endpoint` (Optional) - Webhook endpoint URL.
- `component_ids` (Optional) - List of component IDs to subscribe to.
- `skip_confirmation_notification` (Optional) - Skip sending a confirmation notification.

## Attribute Reference

- `id` - Subscriber identifier.
- `mode` - Subscriber mode. One of `email`, `sms`, `slack`, `webhook`, `microsoft_teams`.

## Import

```shell
terraform import statuspage_subscriber.example page_id/subscriber_id
```
