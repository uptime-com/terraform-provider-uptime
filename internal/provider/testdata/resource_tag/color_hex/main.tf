variable tag {
  type = string
}

variable color_hex {
  type = string
}

resource uptime_tag test {
  tag       = var.tag
  color_hex = var.color_hex
}
