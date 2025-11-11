variable device_name {
  type = string
}

variable uuid {
  type = string
}

variable app_key {
  type      = string
  sensitive = true
}

resource uptime_push_notification test {
  device_name = var.device_name
  uuid        = var.uuid
  app_key     = var.app_key
}

data uptime_push_notifications test {
  depends_on = [uptime_push_notification.test]
}

locals {
  filtered_profiles = [
    for profile in data.uptime_push_notifications.test.profiles :
    profile if profile.device_name == var.device_name
  ]
}

output filtered_count {
  value = length(local.filtered_profiles)
}

output filtered_device_name {
  value = length(local.filtered_profiles) > 0 ? local.filtered_profiles[0].device_name : ""
}

output filtered_uuid {
  value = length(local.filtered_profiles) > 0 ? local.filtered_profiles[0].uuid : ""
}
