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

variable contact_groups {
  type = list(string)
}

resource uptime_check_whois test {
  name          = var.name
  address       = var.address
  expect_string = var.expect_string

  contact_groups = var.contact_groups
}
