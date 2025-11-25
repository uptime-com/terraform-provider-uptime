resource "uptime_integration_status" "example" {
  name          = "My Status.io Integration"
  statuspage_id = "your-statuspage-id"
  api_id        = "your-api-id"
  api_key       = "your-api-key"
  component     = "component-id"
  container     = "container-id"
  metric        = "metric-id"
}
