variable name {
  type = string
}

variable address {
  type = string
}

variable sla_uptime {
  type = string
}

resource uptime_check_rum2 test {
  name       = var.name
  address    = var.address
  sla_uptime = var.sla_uptime
}
