resource "uptime_maintenance_schedule" "test" {
  name                            = "tf-acc-rrule"
  schedule_type                   = "RRULE"
  starts_at                       = "2030-01-01T02:00:00Z"
  rrule                           = "FREQ=WEEKLY;BYDAY=SA"
  duration_minutes                = 120
  is_active                       = true
  pause_checks_during_maintenance = true
}
