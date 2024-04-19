variable name {
  type = string
}

variable script {
  type    = string
  default = <<SCRIPT
[{"step_def": "C_OPEN_URL", "values": {"url": "https://host1"}}]
SCRIPT
}

variable locations {
  type    = list(string)
}

resource uptime_check_transaction test {
  name      = var.name
  script    = var.script
  locations = var.locations
}
