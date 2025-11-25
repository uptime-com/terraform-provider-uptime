resource "uptime_integration_statuspage" "example" {
  name      = "My Statuspage.io Integration"
  api_key   = "your-statuspage-api-key"
  page      = "your-page-id"
  component = "component-id"
  metric    = "metric-id"
}
