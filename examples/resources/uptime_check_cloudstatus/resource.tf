# Legacy: monitor a single named cloud component (deprecated server-side)
resource "uptime_check_cloudstatus" "legacy" {
  name         = "AWS EC2 us-east-1"
  service_name = "aws-ec2-us-east-1"
}

# Group monitoring: track every service in a cloud status group
resource "uptime_check_cloudstatus" "group_all" {
  name            = "AWS Status (all services)"
  group           = 12
  monitoring_type = "ALL"
}

# Group monitoring: track only specific services from a group
resource "uptime_check_cloudstatus" "group_specific" {
  name                = "AWS critical paths"
  contact_groups      = ["nobody"]
  group               = 12
  monitoring_type     = "SPECIFIC"
  services            = [101, 102]
  service_titles      = ["AWS EC2", "AWS RDS"]
  notify_only_on_down = true
}
