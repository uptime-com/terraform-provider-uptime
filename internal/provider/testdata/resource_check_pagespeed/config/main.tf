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


variable "pagespeed_config_exclude_urls" {
  type = string
}

resource "uptime_check_pagespeed" "test" {
  name   = var.name
  script = var.script
  config = {
    exclude_urls = var.pagespeed_config_exclude_urls
  }
}
