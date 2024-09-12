variable "name" {
  type = string
}

variable "sla_report_name" {
  type = string
}

variable "recipient_emails" {
  type = set(string)
}

resource "uptime_sla_report" "test" {
  name = var.sla_report_name
}

resource "uptime_scheduled_report" "test" {
  depends_on       = [uptime_sla_report.test]
  name             = var.name
  sla_report       = var.sla_report_name
  recipient_emails = var.recipient_emails
}
