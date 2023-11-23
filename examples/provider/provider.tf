terraform {
  required_providers {
    uptime = {
      source  = "uptime-com/uptime"
      version = "~> 2.1"
    }
  }
}

variable "uptime_token" {
  type = string
}

provider "uptime" {
  token = var.uptime_token
}

data uptime_locations all {}

resource random_integer location {
  min = 3
  max = length(data.uptime_locations.all.locations) - 1
}
#
resource "uptime_check_http" "http" {
  name           = "Example HTTP Check"
  address        = "https://example.com"
  contact_groups = ["Default"]

  interval  = 5
  locations = [data.uptime_locations.all.locations.*.location[random_integer.location.result]]
}
