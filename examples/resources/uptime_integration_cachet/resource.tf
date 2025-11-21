resource "uptime_integration_cachet" "example" {
  name       = "My Cachet Integration"
  cachet_url = "https://status.example.com"
  token      = "your-cachet-api-token"
  component  = "1"
  metric     = "2"
}
