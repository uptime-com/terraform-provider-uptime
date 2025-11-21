variable name {
  type = string
}

variable api_key {
  type = string
}

variable dataset_name {
  type = string
}

resource uptime_integration_geckoboard test {
  name         = var.name
  api_key      = var.api_key
  dataset_name = var.dataset_name
}
