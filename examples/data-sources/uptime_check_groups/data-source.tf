# Retrieve all check groups
data "uptime_check_groups" "all" {}

# Output all check groups
output "all_groups" {
  value = data.uptime_check_groups.all.check_groups
}

# Filter check groups by name pattern
locals {
  production_groups = [
    for group in data.uptime_check_groups.all.check_groups :
    group if can(regex("production", lower(group.name)))
  ]
}

output "production_groups" {
  value       = local.production_groups
  description = "Check groups related to production"
}
