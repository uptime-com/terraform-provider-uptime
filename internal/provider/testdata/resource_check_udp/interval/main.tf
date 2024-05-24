variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable port {
  type = number
}

variable send_string {
  type    = string
  default = "ping"
}

variable expect_string {
  type    = string
  default = "pong"
}

variable interval {
  type    = number
  default = 10
}

resource uptime_check_udp test {
  name           = var.name
  address        = var.address
  port           = var.port
  send_string    = var.send_string
  expect_string  = var.expect_string
  interval       = var.interval
}
