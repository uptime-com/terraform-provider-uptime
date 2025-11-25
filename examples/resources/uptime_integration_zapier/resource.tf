resource "uptime_integration_zapier" "example" {
  name        = "My Zapier Integration"
  webhook_url = "https://hooks.zapier.com/hooks/catch/your-webhook-url"
}
