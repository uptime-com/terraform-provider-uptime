variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable expect_string {
  type    = string
  default = "example.com"
}

resource uptime_check_whois test {
  name          = var.name
  address       = var.address
  expect_string = var.expect_string
}
