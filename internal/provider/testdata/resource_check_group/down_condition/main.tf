variable "name" {
  type = string
}

variable "down_condition" {
  type = string
}

resource "uptime_check_group" "test" {
  name   = var.name
  config = {
    down_condition = var.down_condition
  }
  contact_groups = []
}
