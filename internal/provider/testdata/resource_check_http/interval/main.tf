variable name {
  type = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable interval {
  type = number
}

resource uptime_check_http test {
  name     = var.name
  address  = var.address
  interval = var.interval
}
