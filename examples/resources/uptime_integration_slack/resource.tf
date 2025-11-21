resource "uptime_integration_slack" "example" {
  name        = "My Slack Integration"
  webhook_url = "https://hooks.slack.com/services/your-webhook-url"
  channel     = "#monitoring"
}
