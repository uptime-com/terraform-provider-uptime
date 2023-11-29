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

variable "metrics_show_section" {
  type = bool
}

variable "metrics_for_all_checks" {
  type = bool
}

resource "uptime_check_api" "metrics" {
  name   = var.check_name
  script = var.script
}

resource "uptime_dashboard" "metrics" {
  depends_on = [uptime_check_api.metrics]
  name       = var.name
  metrics = {
    show_section    = var.metrics_show_section
    for_all_checks  = var.metrics_for_all_checks
  }
  alerts     = {}
  services = {
    show = {}
    sort = {}
  }
  selected = {
    services = [uptime_check_api.metrics.name]
  }
}
