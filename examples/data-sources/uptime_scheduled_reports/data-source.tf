# Retrieve all scheduled reports configured in your Uptime.com account
data "uptime_scheduled_reports" "all" {}

# Filter for enabled weekly reports
locals {
  weekly_reports = [
    for report in data.uptime_scheduled_reports.all.scheduled_reports :
    report if report.recurrence == "weekly" && report.is_enabled
  ]
}

# Filter for a specific report by name
locals {
  monthly_sla_report = [
    for report in data.uptime_scheduled_reports.all.scheduled_reports :
    report if can(regex("monthly.*sla", lower(report.name)))
  ][0]
}

# Output report information
output "active_weekly_reports" {
  value       = length(local.weekly_reports)
  description = "Number of active weekly scheduled reports"
}

output "monthly_report_recipients" {
  value       = local.monthly_sla_report.recipient_emails
  description = "Email recipients for monthly SLA report"
}
