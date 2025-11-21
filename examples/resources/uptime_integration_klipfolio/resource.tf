resource "uptime_integration_klipfolio" "example" {
  name             = "My Klipfolio Integration"
  api_key          = "your-klipfolio-api-key"
  data_source_name = "uptime-monitoring"
}
