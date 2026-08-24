variable name {
  type = string
}

variable script {
  type    = string
  default = <<SCRIPT
[
  {
    "step_def": "C_GET",
    "values": {
      "url": "https://example.com/"
    }
  }
]
SCRIPT
}

variable use_ip_version {
  type = string
}

resource uptime_check_api test {
  name           = var.name
  script         = var.script
  use_ip_version = var.use_ip_version
}
