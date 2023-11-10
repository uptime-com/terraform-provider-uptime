variable name {
  type = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable sla_uptime {
  type = number
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
  sla     = {
    uptime = var.sla_uptime
  }
}
