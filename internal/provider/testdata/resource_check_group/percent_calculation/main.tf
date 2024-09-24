variable "name" {
  type = string
}

variable "uptime_percent_calculation" {
  type = string
}

resource "uptime_check_group" "test" {
  name   = var.name
  config = {
    uptime_percent_calculation = var.uptime_percent_calculation
  }
  contact_groups = []
}
