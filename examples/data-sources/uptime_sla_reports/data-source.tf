# Retrieve all SLA reports configured in your Uptime.com account
data "uptime_sla_reports" "all" {}

# Filter for reports showing uptime section
locals {
  uptime_reports = [
    for report in data.uptime_sla_reports.all.sla_reports :
    report if report.show_uptime_section
  ]
}

# Filter for a specific report by name
locals {
  production_sla = [
    for report in data.uptime_sla_reports.all.sla_reports :
    report if report.name == "Production Services SLA"
  ][0]
}

# Use an SLA report in a scheduled report
resource "uptime_scheduled_report" "weekly" {
  name            = "Weekly Production SLA"
  sla_report      = local.production_sla.url
  recurrence      = "weekly"
  file_type       = "pdf"
  recipient_emails = ["team@example.com"]
}

# Output report information
output "production_sla_url" {
  value       = local.production_sla.stats_url
  description = "Statistics URL for production SLA report"
}
