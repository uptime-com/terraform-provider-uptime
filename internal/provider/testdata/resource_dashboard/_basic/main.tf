variable "name" {
  type = string
}

variable "check_name" {
  type = string
}

variable "script" {
  type    = string
  default = <<SCRIPT
[
  {
    "step_def": "C_GET",
    "values": {
      "url": "https://example.com/"
    }
  }
]
SCRIPT
}

resource "uptime_check_api" "test" {
  name   = var.check_name
  script = var.script
}

resource "uptime_dashboard" "basic" {
  depends_on = [uptime_check_api.test]
  name       = var.name
  alerts     = {}
  services = {
    show = {}
    sort = {}
  }
  selected = {
    services = [uptime_check_api.test.name]
  }
}
