variable name {
  type = string
}

variable cachet_url {
  type = string
}

variable token {
  type = string
}

variable component {
  type = string
}

resource uptime_integration_cachet test {
  name       = var.name
  cachet_url = var.cachet_url
  token      = var.token
  component  = var.component
}
