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

variable contact_groups {
  type = list(string)
}

resource uptime_check_udp test {
  name           = var.name
  address        = var.address
  port           = var.port
  send_string    = var.send_string
  expect_string  = var.expect_string
  contact_groups = var.contact_groups
}
