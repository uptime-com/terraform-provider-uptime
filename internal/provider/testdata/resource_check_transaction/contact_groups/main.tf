variable name {
  type = string
}

variable script {
  type    = string
  default = <<SCRIPT
[{"step_def": "C_OPEN_URL", "values": {"url": "https://host1"}}]
SCRIPT
}

variable contact_groups {
  type = list(string)
}

resource uptime_check_transaction test {
  name           = var.name
  script         = var.script
  contact_groups = var.contact_groups
}
