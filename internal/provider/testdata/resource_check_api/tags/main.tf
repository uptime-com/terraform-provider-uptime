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

variable tags_create {
  type = list(string)
}

resource uptime_tag test {
  count     = length(var.tags_create)
  tag       = var.tags_create[count.index]
  color_hex = "#000000"
}

variable tags_use {
  type = list(string)
}

resource uptime_check_api test {
  name   = var.name
  script = var.script
  tags   = var.tags_use
}
