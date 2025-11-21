variable name {
  type = string
}

variable wavefront_url {
  type = string
}

variable api_token {
  type = string
}

resource uptime_integration_wavefront test {
  name          = var.name
  wavefront_url = var.wavefront_url
  api_token     = var.api_token
}
