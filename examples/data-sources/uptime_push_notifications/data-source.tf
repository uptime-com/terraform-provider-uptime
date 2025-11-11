# Retrieve all push notification profiles
data "uptime_push_notifications" "all" {}

# Output the total number of registered devices
output "total_devices" {
  value = length(data.uptime_push_notifications.all.profiles)
}

# Filter for a specific device by name
locals {
  my_iphone = [
    for profile in data.uptime_push_notifications.all.profiles :
    profile if profile.device_name == "My iPhone"
  ]
}

# Output the filtered device details (if found)
output "my_iphone_id" {
  value = length(local.my_iphone) > 0 ? local.my_iphone[0].id : null
}

output "my_iphone_uuid" {
  value = length(local.my_iphone) > 0 ? local.my_iphone[0].uuid : null
}

# List all device names
output "all_device_names" {
  value = [for profile in data.uptime_push_notifications.all.profiles : profile.device_name]
}

# Find devices in specific contact groups
locals {
  production_devices = [
    for profile in data.uptime_push_notifications.all.profiles :
    profile if contains(profile.contact_groups, "Production")
  ]
}

output "production_device_count" {
  value = length(local.production_devices)
}
