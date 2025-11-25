resource "uptime_integration_datadog" "example" {
  name    = "My Datadog Integration"
  api_key = "your-datadog-api-key"
  app_key = "your-datadog-app-key"
  region  = "us"
}
