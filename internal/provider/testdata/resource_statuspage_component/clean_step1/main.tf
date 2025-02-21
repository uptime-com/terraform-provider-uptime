variable "name" {
  type = string
}

variable "component_name" {
  type = string
}

resource "uptime_check_api" "test" {
  name   = "test"
  script = <<SCRIPT
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

resource "uptime_statuspage_component" "group" {
  statuspage_id = uptime_statuspage.test.id
  name          = "group component"
  is_group      = true
}

resource "uptime_statuspage" "test" {
  name = var.name
}

resource "uptime_statuspage_component" "test" {
  statuspage_id = uptime_statuspage.test.id
  name          = var.component_name
  description   = "Test Component"
  service_id    = uptime_check_api.test.id
  group_id      = uptime_statuspage_component.group.id
}
