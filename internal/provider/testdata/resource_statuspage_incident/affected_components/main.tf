variable "name" {
  type = string
}

variable "check_name" {
  type = string
}

variable "incident_name" {
  type = string
}

variable "incident_state" {
  type = string
}

variable "incident_component_status" {
  type = string
}

variable "starts_at" {
  type = string
}

variable "updated_at" {
  type    = string
  default = null
}

resource "uptime_check_api" "test" {
  name   = var.check_name
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

resource "uptime_statuspage" "test" {
  name = var.name
}

resource "uptime_statuspage_component" "test" {
  statuspage_id = uptime_statuspage.test.id
  name          = var.check_name
  service_id    = uptime_check_api.test.id
}

resource "uptime_statuspage_incident" "test" {
  statuspage_id = uptime_statuspage.test.id
  name          = var.incident_name
  incident_type = "INCIDENT"
  starts_at     = var.starts_at
  updates = [
    {
      updated_at     = var.updated_at
      incident_state = var.incident_state
    }
  ]
  affected_components = [
    {
      status       = var.incident_component_status
      component_id = uptime_statuspage_component.test.id
    }
  ]
}
