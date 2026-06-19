resource "uptime_contact" "test" {
  name       = "tf-acc-maint-contact"
  email_list = ["tf-acc-maint@example.com"]
}

resource "uptime_maintenance_schedule" "test" {
  name             = "tf-acc-maint-sched"
  schedule_type    = "ONE_OFF"
  starts_at        = "2030-01-01T02:00:00Z"
  duration_minutes = 60
}

resource "uptime_maintenance_notification" "test" {
  schedule_id    = uptime_maintenance_schedule.test.id
  offset         = -1800
  event          = "START"
  contact_groups = [uptime_contact.test.id]
}
