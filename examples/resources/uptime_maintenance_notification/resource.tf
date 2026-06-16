resource "uptime_maintenance_notification" "before_start" {
  schedule_id    = uptime_maintenance_schedule.weekly_patching.id
  offset         = -1800 # 30 minutes before
  event          = "START"
  contact_groups = [uptime_contact.oncall.id]
}
