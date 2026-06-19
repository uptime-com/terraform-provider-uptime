resource "uptime_maintenance_schedule" "weekly_patching" {
  name             = "Weekly Patching Window"
  schedule_type    = "RRULE"
  starts_at        = "2026-06-20T02:00:00Z"
  rrule            = "FREQ=WEEKLY;BYDAY=SA"
  duration_minutes = 120

  pause_checks_during_maintenance = true

  services = [uptime_check_http.api.id]
  tags     = [uptime_tag.production.id]
}
