variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

resource uptime_check_icmp test {
  name    = var.name
  address = var.address
}
