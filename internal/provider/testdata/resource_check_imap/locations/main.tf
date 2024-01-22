variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable locations {
  type = list(string)
}
resource uptime_check_imap test {
  name      = var.name
  address   = var.address
  locations = var.locations
}
