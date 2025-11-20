# Retrieve all users in your Uptime.com account
data "uptime_users" "all" {}

# Filter for admin users
locals {
  admin_users = [
    for user in data.uptime_users.all.users :
    user if user.access_level == "admin" && user.is_active
  ]
}

# Filter for a specific user by email
locals {
  ops_user = [
    for user in data.uptime_users.all.users :
    user if user.email == "ops@example.com"
  ][0]
}

# Output user information
output "admin_user_count" {
  value       = length(local.admin_users)
  description = "Number of active admin users"
}

output "ops_user_id" {
  value       = local.ops_user.id
  description = "User ID for ops team member"
}
