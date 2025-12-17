# First create an SLA report that the scheduled report will send
resource "uptime_sla_report" "production" {
  name = "Production SLA Report"
}

# Basic scheduled report (weekly)
resource "uptime_scheduled_report" "example" {
  name             = "Weekly Status Report"
  sla_report       = uptime_sla_report.production.name
  recurrence       = "WEEKLY"
  recipient_emails = ["team@example.com"]
}

# Daily report
resource "uptime_scheduled_report" "daily" {
  name             = "Daily Operations Report"
  sla_report       = uptime_sla_report.production.name
  recurrence       = "DAILY"
  recipient_emails = ["ops@example.com"]
}

# Monthly report with PDF format
resource "uptime_scheduled_report" "monthly" {
  name             = "Monthly Executive Summary"
  sla_report       = uptime_sla_report.production.name
  recurrence       = "MONTHLY"
  recipient_emails = ["executives@example.com"]
  file_type        = "PDF"
}
