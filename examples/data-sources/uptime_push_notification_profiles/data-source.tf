# Retrieve all push notification profiles (mobile devices) in your account
data "uptime_push_notification_profiles" "all" {}

# Filter for a specific device by name
locals {
  my_phone = [
    for profile in data.uptime_push_notification_profiles.all.push_notification_profiles :
    profile if can(regex("iphone", lower(profile.device_name)))
  ][0]
}

# Filter for devices by user
locals {
  admin_devices = [
    for profile in data.uptime_push_notification_profiles.all.push_notification_profiles :
    profile if can(regex("admin", lower(profile.display_name)))
  ]
}

# Output device information
output "my_device_id" {
  value       = local.my_phone.id
  description = "Push notification profile ID for my iPhone"
}

output "admin_device_count" {
  value       = length(local.admin_devices)
  description = "Number of mobile devices registered to admins"
}

output "all_device_names" {
  value = [
    for profile in data.uptime_push_notification_profiles.all.push_notification_profiles :
    "${profile.display_name} (${profile.device_name})"
  ]
  description = "List of all registered mobile devices"
}
