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

resource "uptime_check_pagespeed" "test" {
  name   = var.name
  script = var.script
}
