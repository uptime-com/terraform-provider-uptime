variable "name" {
  type = string
}

variable "email" {
  type = string
}

variable "first_name" {
  type = string
}

variable "last_name" {
  type = string
}

variable "is_active" {
  type = bool
}

resource "uptime_statuspage" "test" {
  name = var.name
}

resource "uptime_statuspage_user" "test" {
  statuspage_id = uptime_statuspage.test.id
  email         = var.email
  first_name    = var.first_name
  last_name     = var.last_name
  is_active     = var.is_active
}
