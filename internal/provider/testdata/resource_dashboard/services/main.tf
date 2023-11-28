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

resource "uptime_check_api" "services" {
  name   = var.check_name
  script = var.script
}

variable "services_show_section" {
  type = string
}

variable "services_num_to_show" {
  type = number
}

variable "services_include_up" {
  type = bool
}

variable "services_include_down" {
  type = bool
}

variable "services_include_paused" {
  type = bool
}

variable "services_include_maintenance" {
  type = bool
}

variable "services_primary_sort" {
  type = string
}

variable "services_secondary_sort" {
  type = string
}

variable "services_show_uptime" {
  type = bool
}

variable "services_show_response_time" {
  type = bool
}

resource "uptime_dashboard" "services" {
  depends_on = [uptime_check_api.services]
  name       = var.name
  alerts     = {}
  services = {
    show_section = var.services_show_section
    num_to_show  = var.services_num_to_show
    include = {
      up          = var.services_include_up
      down        = var.services_include_down
      paused      = var.services_include_paused
      maintenance = var.services_include_maintenance
    }
    sort = {
      primary   = var.services_primary_sort
      secondary = var.services_secondary_sort
    }
    show = {
      uptime        = var.services_show_uptime
      response_time = var.services_show_response_time
    }
  }
  selected = {
    services = [uptime_check_api.services.name]
  }
}
