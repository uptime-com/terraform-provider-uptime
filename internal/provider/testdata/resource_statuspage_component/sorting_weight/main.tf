variable "name" {
  type = string
}

variable "component_name" {
  type = string
}

variable "sorting_weight" {
  type = number
}

resource "uptime_statuspage" "test" {
  name = var.name
}

resource "uptime_statuspage_component" "test" {
  statuspage_id  = uptime_statuspage.test.id
  name           = var.component_name
  sorting_weight = var.sorting_weight
}
