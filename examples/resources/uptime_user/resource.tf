# Create a basic user
resource "uptime_user" "example" {
  email      = "user@example.com"
  first_name = "John"
  last_name  = "Doe"
  password   = "secure-password"
}

# Create an admin user with API access
resource "uptime_user" "admin" {
  email                = "admin@example.com"
  first_name           = "Jane"
  last_name            = "Admin"
  password             = "secure-password"
  access_level         = "admin"
  is_api_enabled       = true
  notify_paid_invoices = true
}

# Create a user with subaccount access
resource "uptime_user" "subaccount_user" {
  email      = "subuser@example.com"
  first_name = "Bob"
  last_name  = "Smith"
  password   = "secure-password"
  assigned_subaccounts = [
    "https://uptime.com/api/v1/subaccounts/123/",
    "https://uptime.com/api/v1/subaccounts/456/",
  ]
}
