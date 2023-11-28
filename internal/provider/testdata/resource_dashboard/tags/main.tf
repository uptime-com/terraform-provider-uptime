variable "name" {
  type = string
}

variable "tags_create" {
  type = list(string)
}

resource "uptime_tag" "tags" {
  count     = length(var.tags_create)
  tag       = var.tags_create[count.index]
  color_hex = "#000000"
}

variable "tags_use" {
  type = list(string)
}

resource "uptime_dashboard" "tags" {
  depends_on = [uptime_tag.tags]
  name       = var.name
  alerts     = {}
  services = {
    show = {}
    sort = {}
  }
  selected = {
    tags = var.tags_use
  }
}
