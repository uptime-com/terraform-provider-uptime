variable name {
  type = string
}

variable address {
  type = string
}

variable contact_groups {
  type = list(string)
}

resource uptime_check_rum2 test {
  name           = var.name
  address        = var.address
  contact_groups = var.contact_groups
}
