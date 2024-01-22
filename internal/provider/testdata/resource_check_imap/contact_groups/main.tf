variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable contact_groups {
  type = list(string)
}

resource uptime_check_imap test {
  name           = var.name
  address        = var.address
  contact_groups = var.contact_groups
}
