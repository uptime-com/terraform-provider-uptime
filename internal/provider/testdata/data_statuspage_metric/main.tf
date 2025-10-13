variable statuspage_name {
  type = string
}

variable metric_name {
  type = string
}

variable check_name {
  type = string
}

resource uptime_statuspage test {
  name = var.statuspage_name
}

resource uptime_check_http test {
  name    = var.check_name
  address = "https://example.com"
}

resource uptime_statuspage_metric test {
  statuspage_id = uptime_statuspage.test.id
  name          = var.metric_name
  service_id    = uptime_check_http.test.id
  is_visible    = true
}

data uptime_statuspage_metrics test {
  statuspage_id = uptime_statuspage.test.id
  depends_on    = [uptime_statuspage_metric.test]
}

locals {
  filtered_metrics = [
    for metric in data.uptime_statuspage_metrics.test.metrics :
    metric if metric.name == var.metric_name
  ]
}

output filtered_count {
  value = length(local.filtered_metrics)
}

output filtered_metric_name {
  value = length(local.filtered_metrics) > 0 ? local.filtered_metrics[0].name : ""
}
