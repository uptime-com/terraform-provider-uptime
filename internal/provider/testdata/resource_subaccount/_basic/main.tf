variable "name" {
  type = string
}

resource "uptime_subaccount" "test" {
  name = var.name
}
