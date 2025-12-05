variable "name" {
  type = string
}

resource "uptime_check_group" "test" {
  name   = var.name
  config = {}
}

data "uptime_check_groups" "test" {
  depends_on = [uptime_check_group.test]
}

locals {
  filtered_check_groups = [
    for check_group in data.uptime_check_groups.test.check_groups :
    check_group if check_group.name == var.name
  ]
}

output "filtered_count" {
  value = tostring(length(local.filtered_check_groups))
}

output "filtered_check_group_name" {
  value = length(local.filtered_check_groups) > 0 ? local.filtered_check_groups[0].name : ""
}
