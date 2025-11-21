resource "uptime_integration_webhook" "example" {
  name                = "My Webhook Integration"
  postback_url        = "https://example.com/webhook"
  headers             = jsonencode({ "Authorization" = "Bearer token" })
  use_legacy_payload  = false
}
