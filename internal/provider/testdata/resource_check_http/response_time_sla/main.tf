variable name {
  type = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable response_time_sla {
  type = string
}

resource uptime_check_http test {
  name              = var.name
  address           = var.address
  response_time_sla = var.response_time_sla
}
