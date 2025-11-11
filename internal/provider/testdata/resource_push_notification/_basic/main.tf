variable device_name {
  type = string
}

variable uuid {
  type = string
}

variable app_key {
  type = string
}

resource uptime_push_notification test {
  device_name = var.device_name
  uuid        = var.uuid
  app_key     = var.app_key
}
