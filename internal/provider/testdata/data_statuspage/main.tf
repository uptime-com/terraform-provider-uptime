variable name {
  type = string
}

resource uptime_statuspage test {
  name = var.name
}

data uptime_statuspages test {
  depends_on = [uptime_statuspage.test]
}

locals {
  filtered_statuspages = [
    for sp in data.uptime_statuspages.test.statuspages :
    sp if sp.name == var.name
  ]
}

output filtered_count {
  value = length(local.filtered_statuspages)
}

output filtered_statuspage_name {
  value = length(local.filtered_statuspages) > 0 ? local.filtered_statuspages[0].name : ""
}
