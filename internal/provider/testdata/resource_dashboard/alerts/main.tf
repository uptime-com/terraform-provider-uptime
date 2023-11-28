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

resource "uptime_check_api" "alerts" {
  name   = var.check_name
  script = var.script
}

variable "alerts_show_section" {
  type = bool
}

variable "alerts_for_all_checks" {
  type = bool
}

variable "alerts_num_to_show" {
  type = number
}

variable "alerts_include_ignored" {
  type = bool
}

variable "alerts_include_resolved" {
  type = bool
}

resource "uptime_dashboard" "alerts" {
  depends_on = [uptime_check_api.alerts]
  name       = var.name
  alerts = {
    show           = var.alerts_show_section
    for_all_checks = var.alerts_for_all_checks
    num_to_show    = var.alerts_num_to_show
    include = {
      ignored  = var.alerts_include_ignored
      resolved = var.alerts_include_resolved
    }
  }
  services = {
    sort = {}
    show = {}
  }
  selected = {
    services = [uptime_check_api.alerts.name]
  }
}
