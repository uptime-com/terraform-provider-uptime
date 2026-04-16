variable name {
  type = string
}

variable service_name {
  type    = string
  default = "Amazon Service"
}

resource uptime_check_cloudstatus test {
  name         = var.name
  service_name = var.service_name
}
