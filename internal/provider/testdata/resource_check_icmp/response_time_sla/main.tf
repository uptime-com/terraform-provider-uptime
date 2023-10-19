variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable response_time_sla {
  type = string
}

resource uptime_check_icmp test {
  name              = var.name
  address           = var.address
  response_time_sla = var.response_time_sla
}
