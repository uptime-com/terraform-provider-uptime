resource uptime_check_blacklist test {
  name           = var.name
  address        = var.address
}

variable name {
  type = string
}

variable address {
  type = string
}
