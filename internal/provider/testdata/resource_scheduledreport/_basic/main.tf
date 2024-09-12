variable "name" {
  type = string
}

variable "sla_report_name" {
  type = string
}

resource "uptime_sla_report" "test" {
  name = var.sla_report_name
}

resource "uptime_scheduled_report" "test" {
  depends_on = [uptime_sla_report.test]
  name       = var.name
  sla_report = var.sla_report_name
}
