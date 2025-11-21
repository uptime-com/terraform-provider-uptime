variable name {
  type = string
}

variable api_key {
  type = string
}

variable data_source_name {
  type = string
}

resource uptime_integration_klipfolio test {
  name             = var.name
  api_key          = var.api_key
  data_source_name = var.data_source_name
}
