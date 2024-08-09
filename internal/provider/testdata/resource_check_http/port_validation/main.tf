variable name {
  type    = string
}

variable address {
  type    = string
  default = "https://example.com:8383"
}

variable port {
  type    = number
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
  port    = var.port
}

