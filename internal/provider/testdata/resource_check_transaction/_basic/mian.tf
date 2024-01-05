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

resource uptime_check_transaction test {
  name   = var.name
  script = var.script
}
