variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable interval {
  type    = number
  default = 10
}

resource uptime_check_icmp test {
  name     = var.name
  address  = var.address
  interval = var.interval
}
