# Basic transaction check (browser-based synthetic monitoring)
resource "uptime_check_transaction" "example" {
  name = "Login Flow Check"
  script = jsonencode([
    {
      step_def = "C_OPEN_URL"
      values = {
        url = "https://app.example.com/login"
      }
    }
  ])
}

# Transaction check with multiple URL steps
resource "uptime_check_transaction" "multi_page" {
  name = "Multi-Page Flow"
  script = jsonencode([
    {
      step_def = "C_OPEN_URL"
      values = {
        url = "https://app.example.com/"
      }
    },
    {
      step_def = "C_OPEN_URL"
      values = {
        url = "https://app.example.com/dashboard"
      }
    },
    {
      step_def = "C_OPEN_URL"
      values = {
        url = "https://app.example.com/profile"
      }
    }
  ])
}

# Transaction check with full configuration
resource "uptime_check_transaction" "full" {
  name = "Critical User Journey"
  script = jsonencode([
    {
      step_def = "C_OPEN_URL"
      values = {
        url = "https://shop.example.com"
      }
    },
    {
      step_def = "C_OPEN_URL"
      values = {
        url = "https://shop.example.com/products"
      }
    },
    {
      step_def = "C_OPEN_URL"
      values = {
        url = "https://shop.example.com/cart"
      }
    }
  ])
  interval       = 15
  contact_groups = ["nobody"]
  sla = {
    uptime  = "0.999"
    latency = "5s"
  }
}
