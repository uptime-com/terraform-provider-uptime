resource "uptime_integration_victorops" "example" {
  name        = "My VictorOps Integration"
  service_key = "your-victorops-service-key"
  routing_key = "production"
}
