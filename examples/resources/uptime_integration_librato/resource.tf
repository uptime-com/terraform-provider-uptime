resource "uptime_integration_librato" "example" {
  name        = "My Librato Integration"
  email       = "user@example.com"
  api_token   = "your-librato-api-token"
  metric_name = "uptime.response_time"
}
