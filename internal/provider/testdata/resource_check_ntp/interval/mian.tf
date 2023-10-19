variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable interval {
  type = number
}

resource uptime_check_ntp test {
  name     = var.name
  address  = var.address
  interval = var.interval
}
