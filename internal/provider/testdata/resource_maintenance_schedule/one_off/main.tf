resource "uptime_maintenance_schedule" "test" {
  name             = "tf-acc-one-off"
  schedule_type    = "ONE_OFF"
  starts_at        = "2030-01-01T02:00:00Z"
  duration_minutes = 120
}
