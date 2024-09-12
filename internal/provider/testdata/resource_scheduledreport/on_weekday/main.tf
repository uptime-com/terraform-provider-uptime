variable "name" {
  type = string
}

variable "sla_report_name" {
  type = string
}

variable "on_weekday" {
  type = number
}

resource "uptime_sla_report" "test" {
  name = var.sla_report_name
}

resource "uptime_scheduled_report" "test" {
  depends_on = [uptime_sla_report.test]
  name       = var.name
  sla_report = var.sla_report_name
  on_weekday = var.on_weekday
}
