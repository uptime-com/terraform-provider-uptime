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

variable "is_pinned" {
  type = bool
}

variable "ordering" {
  type = number
}

resource "uptime_check_api" "root" {
  name   = var.check_name
  script = var.script
}

resource "uptime_dashboard" "root" {
  depends_on = [uptime_check_api.root]
  name       = var.name
  ordering   = var.ordering
  is_pinned  = var.is_pinned
  alerts     = {}
  services = {
    show = {}
    sort = {}
  }
  selected = {
    services = [uptime_check_api.root.name]
  }
}
