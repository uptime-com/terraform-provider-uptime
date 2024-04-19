variable name {
  type = string
}

variable script {
  type    = string
  default = <<SCRIPT
[{"step_def": "C_OPEN_URL", "values": {"url": "https://host1"}}]
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

resource uptime_check_transaction test {
  depends_on = [uptime_tag.test]
  name       = var.name
  script     = var.script
  tags       = var.tags_use
}
