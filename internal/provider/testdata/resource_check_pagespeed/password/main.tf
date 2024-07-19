variable "name" {
  type = string
}

variable "script" {
  type    = string
  default = <<SCRIPT
[{
  "step_def": "C_PAGESPEED_NAVIGATE",
  "values": {
    "url": "https://example.com"
  }
}]
SCRIPT
}

variable password {
  type = string
  sensitive = false
}

resource "uptime_check_pagespeed" "test" {
  name     = var.name
  script   = var.script
  password = var.password
}
