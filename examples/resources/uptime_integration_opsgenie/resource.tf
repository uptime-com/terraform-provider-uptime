resource "uptime_integration_opsgenie" "example" {
  name         = "My OpsGenie Integration"
  api_endpoint = "https://api.opsgenie.com/v1/json/uptime"
  api_key      = "your-opsgenie-api-key"
  teams        = "team1,team2"
  tags         = "monitoring,production"
  auto_resolve = true
}
