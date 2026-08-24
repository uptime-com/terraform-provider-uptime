variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable use_ip_version {
  type = string
}

resource uptime_check_icmp test {
  name           = var.name
  address        = var.address
  use_ip_version = var.use_ip_version
}
