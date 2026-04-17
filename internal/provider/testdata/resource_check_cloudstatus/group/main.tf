variable name {
  type = string
}

variable group {
  type = number
}

variable monitoring_type {
  type    = string
  default = "ALL"
}

variable services {
  type    = set(number)
  default = []
}

variable service_titles {
  type    = set(string)
  default = []
}

resource uptime_check_cloudstatus test {
  name            = var.name
  group           = var.group
  monitoring_type = var.monitoring_type
  services        = var.services
  service_titles  = var.service_titles
}
