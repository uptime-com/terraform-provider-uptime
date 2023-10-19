variable tag {
  type = string
}

variable color_hex {
  type = string
}

resource uptime_tag color_hex {
  tag       = var.tag
  color_hex = var.color_hex
}
