variable name {
  type = string
}

variable address {
  type = string
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

resource uptime_check_rum2 test {
  depends_on = [uptime_tag.test]
  name       = var.name
  address    = var.address
  tags       = var.tags_use
}
