# Basic webhook integration
resource "uptime_integration_webhook" "example" {
  name               = "My Webhook Integration"
  postback_url       = "https://example.com/webhook"
  use_legacy_payload = false
}

# Webhook with custom headers
# Note: headers use newline-delimited "key: value" format, NOT JSON
resource "uptime_integration_webhook" "with_headers" {
  name               = "Authenticated Webhook"
  postback_url       = "https://api.example.com/alerts"
  headers            = "Authorization: Bearer your-token-here"
  use_legacy_payload = false
}
