variable name {
  type = string
}

variable postback_url {
  type = string
}

variable headers {
  type = string
}

resource uptime_integration_webhook test {
  name         = var.name
  postback_url = var.postback_url
  headers      = var.headers
}
