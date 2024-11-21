variable "name" {
  type = string
}

variable "incident_name" {
  type = string
}

variable "incident_state" {
  type = string
}

variable "ends_at" {
  type    = string
  default = null
}

resource "uptime_statuspage" "test" {
  name = var.name
}

resource "uptime_statuspage_incident" "test" {
  statuspage_id  = uptime_statuspage.test.id
  name           = var.incident_name
  incident_type  = "INCIDENT"
  ends_at        = var.ends_at
  updates        = [
    {
      incident_state = var.incident_state
    }
  ]
}
