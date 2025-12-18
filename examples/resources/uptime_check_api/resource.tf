# Basic API check with simple GET request
resource "uptime_check_api" "example" {
  name = "API Health Check"
  script = jsonencode([
    {
      step_def = "C_GET"
      values = {
        url = "https://api.example.com/health"
      }
    }
  ])
}

# API check with multiple steps
resource "uptime_check_api" "multi_step" {
  name = "API Workflow Check"
  script = jsonencode([
    {
      step_def = "C_GET"
      values = {
        url = "https://api.example.com/auth/token"
      }
    },
    {
      step_def = "C_GET"
      values = {
        url = "https://api.example.com/users"
      }
    }
  ])
}

# API check with POST request
resource "uptime_check_api" "post_check" {
  name = "API POST Check"
  script = jsonencode([
    {
      step_def = "C_POST"
      values = {
        url  = "https://api.example.com/data"
        body = "{\"test\": true}"
      }
    }
  ])
}

# API check with full configuration
resource "uptime_check_api" "full" {
  name = "Production API"
  script = jsonencode([
    {
      step_def = "C_GET"
      values = {
        url = "https://api.example.com/status"
      }
    }
  ])
  interval       = 5
  contact_groups = ["nobody"]
  sla = {
    uptime  = "0.999"
    latency = "1s"
  }
}
