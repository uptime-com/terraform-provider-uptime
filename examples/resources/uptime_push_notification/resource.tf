# Configure a push notification profile for receiving mobile alerts
#
# IMPORTANT: To get the required credentials:
# 1. Install the Uptime.com mobile app on your iOS or Android device
# 2. Log in to your Uptime.com account in the app
# 3. Register your device for push notifications
# 4. The app_key and uuid are generated during device registration
#
# Note: The app_key must be exactly 32 characters and is device-specific.
# These credentials authenticate your device to receive push notifications.

resource "uptime_push_notification" "mobile_device" {
  device_name = "My iPhone"
  uuid        = "550e8400-e29b-41d4-a716-446655440000" # Device UUID from mobile app registration
  app_key     = "your-32-character-app-key-here123" # Exactly 32 chars, from mobile app

  # Optional: Specify which contact groups should trigger notifications to this device
  # Defaults to ["Default"] if not specified
  contact_groups = ["Default", "Production"]
}

# Output the profile details (app_key is marked sensitive and won't be shown)
output "push_notification_id" {
  value = uptime_push_notification.mobile_device.id
}

output "push_notification_display_name" {
  value = uptime_push_notification.mobile_device.display_name
}
