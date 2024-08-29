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

variable tags_create {
  type = list(string)
}

resource uptime_tag test {
  count     = length(var.tags_create)
  tag       = var.tags_create[count.index]
  color_hex = "#000000"
}

variable tags_use {
  type = list(string)
}

resource "uptime_check_api" "test" {
  depends_on = [uptime_tag.test]
  name       = var.check_name
  script     = var.script
  tags       = var.tags_use
}

resource "uptime_sla_report" "test" {
  depends_on    = [uptime_check_api.test]
  name          = var.name
  services_tags = var.tags_use
}
