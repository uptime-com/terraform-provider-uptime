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

variable response_time_sla {
  type = string
}

resource uptime_check_api test {
  name              = var.name
  script            = var.script
  response_time_sla = var.response_time_sla
}
