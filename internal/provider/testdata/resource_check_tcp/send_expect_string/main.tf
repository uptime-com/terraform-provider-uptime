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
  type = string
}

variable expect_string {
  type = string
}

resource uptime_check_tcp test {
  name          = var.name
  address       = var.address
  port          = var.port
  send_string   = var.send_string
  expect_string = var.expect_string
}
