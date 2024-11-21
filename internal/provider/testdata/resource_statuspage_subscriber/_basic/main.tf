variable "name" {
  type = string
}

variable "subscriber_target" {
  type = string
}

variable "subscriber_type" {
  type = string
}

resource "uptime_statuspage" "test" {
  name = var.name
}

resource "uptime_statuspage_subscriber" "test" {
  statuspage_id = uptime_statuspage.test.id
  target        = var.subscriber_target
  type          = var.subscriber_type
}
