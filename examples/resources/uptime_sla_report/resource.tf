# Basic SLA report
resource "uptime_sla_report" "example" {
  name = "Production SLA Report"
}

# SLA report with configuration
resource "uptime_sla_report" "configured" {
  name                         = "Configured SLA Report"
  default_date_range           = "LAST_30D"
  show_uptime_section          = true
  show_uptime_sla              = true
  show_response_time_section   = true
  show_response_time_sla       = true
  filter_with_downtime         = true
  filter_uptime_sla_violations = false
}

# SLA report using check from created resources
resource "uptime_check_http" "api" {
  name    = "API Health"
  address = "https://api.example.com/health"
}

resource "uptime_sla_report" "dynamic" {
  name = "Dynamic SLA Report"
  services_selected = [
    { name = uptime_check_http.api.name }
  ]
}
