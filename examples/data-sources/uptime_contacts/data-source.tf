# Retrieve all contacts configured in your Uptime.com account
data "uptime_contacts" "all" {}

# Filter for contacts containing "production" in the name
locals {
  production_contacts = [
    for contact in data.uptime_contacts.all.contacts :
    contact if can(regex("production", lower(contact.name)))
  ]
}

# Use a contact in a check configuration
resource "uptime_check_http" "api" {
  name    = "API Health Check"
  address = "https://api.example.com/health"

  contact_groups = [
    local.production_contacts[0].url
  ]
}

# Output contact details
output "production_contact_emails" {
  value       = local.production_contacts[0].email_list
  description = "Email addresses for production contact group"
}
