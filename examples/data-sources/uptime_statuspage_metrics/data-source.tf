# Reference an existing status page
resource "uptime_statuspage" "main" {
  name = "Production Status"
}

# Retrieve all metrics for the status page
data "uptime_statuspage_metrics" "all" {
  statuspage_id = uptime_statuspage.main.id
}

# Filter for visible metrics only
locals {
  visible_metrics = [
    for metric in data.uptime_statuspage_metrics.all.metrics :
    metric if metric.is_visible
  ]
}

# Output metric names
output "visible_metric_names" {
  value = [for m in local.visible_metrics : m.name]
}
