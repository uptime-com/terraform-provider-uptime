variable "name" {
  type = string
}

variable "metric_name" {
  type = string
}

variable "is_visible" {
  type = bool
}

resource "uptime_check_api" "test" {
  name          = "test"
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

resource "uptime_statuspage_metric" "test" {
  statuspage_id = uptime_statuspage.test.id
  name          = var.metric_name
  service_id    = uptime_check_api.test.id
  is_visible    = var.is_visible
}
