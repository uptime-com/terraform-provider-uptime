variable name {
  type = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable use_ip_version {
  type = string
}

resource uptime_check_http test {
  name           = var.name
  address        = var.address
  use_ip_version = var.use_ip_version
}
