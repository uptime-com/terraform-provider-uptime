resource "uptime_check_group" "test_group" {
  name           = var.group_name
  contact_groups = ["Default"]
  config         = {}
}

resource "uptime_check_http" "test" {
  name           = var.name
  address        = var.address
  contact_groups = []
}
