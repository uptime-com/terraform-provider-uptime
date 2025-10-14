variable statuspage_name {
  type = string
}

variable component_name {
  type = string
}

resource uptime_statuspage test {
  name = var.statuspage_name
}

resource uptime_statuspage_component test {
  statuspage_id = uptime_statuspage.test.id
  name          = var.component_name
  description   = "Test component"
}

data uptime_statuspage_components test {
  statuspage_id = uptime_statuspage.test.id
  depends_on    = [uptime_statuspage_component.test]
}

locals {
  filtered_components = [
    for comp in data.uptime_statuspage_components.test.components :
    comp if comp.name == var.component_name
  ]
}

output filtered_count {
  value = length(local.filtered_components)
}

output filtered_component_name {
  value = length(local.filtered_components) > 0 ? local.filtered_components[0].name : ""
}
