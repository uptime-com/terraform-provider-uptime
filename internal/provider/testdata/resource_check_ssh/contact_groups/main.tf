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

variable contact_groups {
  type = list(string)
}

resource uptime_check_ssh test {
  name           = var.name
  address        = var.address
  port           = var.port
  contact_groups = var.contact_groups
}
