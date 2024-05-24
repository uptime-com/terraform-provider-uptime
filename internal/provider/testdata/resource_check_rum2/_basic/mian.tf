variable name {
  type = string
}

variable address {
  type = string
}

resource uptime_check_rum2 test {
  name    = var.name
  address = var.address
}
