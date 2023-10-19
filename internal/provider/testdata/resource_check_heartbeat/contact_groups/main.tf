variable name {
  type = string
}

variable contact_groups {
  type = list(string)
}

resource uptime_check_heartbeat test {
  name           = var.name
  contact_groups = var.contact_groups
}
