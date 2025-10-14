# Reference an existing status page
resource "uptime_statuspage" "main" {
  name = "Production Status"
}

# Retrieve all users for the status page
data "uptime_statuspage_users" "all" {
  statuspage_id = uptime_statuspage.main.id
}

# Filter for active users only
locals {
  active_users = [
    for user in data.uptime_statuspage_users.all.users :
    user if user.is_active
  ]
}

# Find a specific user by email
locals {
  admin_user = [
    for user in data.uptime_statuspage_users.all.users :
    user if user.email == "admin@example.com"
  ][0]
}

# Output user information
output "active_user_count" {
  value = length(local.active_users)
}

output "admin_user_name" {
  value = "${local.admin_user.first_name} ${local.admin_user.last_name}"
}
