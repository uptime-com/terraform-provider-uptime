resource "uptime_integration_pagerduty" "example" {
  name         = "My PagerDuty Integration"
  service_key  = "your-pagerduty-service-key"
  auto_resolve = true
}
