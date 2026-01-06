# Basic check group
resource "uptime_check_group" "example" {
  name   = "Production Services"
  config = {}
}

# Check group with down condition and contacts
resource "uptime_check_group" "with_condition" {
  name = "API Services"
  config = {
    down_condition = "TWO"
  }
  contact_groups = ["nobody"]
}

# Check group with AVERAGE uptime calculation
# Note: AVERAGE mode cannot have contact groups assigned
# resource "uptime_check_group" "with_calculation" {
#   name = "Critical Infrastructure"
#   config = {
#     uptime_percent_calculation = "AVERAGE"
#   }
# }

# Using check group with HTTP checks (recommended pattern)
# Reference check names from created resources
# Note: services can reference checks by name or ID
# resource "uptime_check_http" "api" {
#   name    = "API Health"
#   address = "https://api.example.com/health"
# }
#
# resource "uptime_check_http" "web" {
#   name    = "Web Frontend"
#   address = "https://www.example.com"
# }
#
# resource "uptime_check_group" "combined" {
#   depends_on = [uptime_check_http.api, uptime_check_http.web]
#   name       = "Web Application"
#   config = {
#     services       = [uptime_check_http.api.name, uptime_check_http.web.name]
#     down_condition = "ALL"
#   }
#   contact_groups = ["nobody"]
# }
