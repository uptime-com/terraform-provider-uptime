variable "name" {
  type = string
}

variable "sla_report_name" {
  type = string
}

variable "at_time" {
  type = number
}

resource "uptime_sla_report" "test" {
  name = var.sla_report_name
}

resource "uptime_scheduled_report" "test" {
  depends_on = [uptime_sla_report.test]
  name       = var.name
  sla_report = var.sla_report_name
  at_time    = var.at_time
}
