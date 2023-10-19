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

variable locations {
  type    = list(string)
}

resource uptime_check_api test {
  name      = var.name
  script    = var.script
  locations = var.locations
}
