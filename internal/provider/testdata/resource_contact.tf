resource uptime_contact test {
  name       = "{{ petname 3 "-" }}"
  email_list = ["example@uptime.com"]
}
// ---
# TODO: Uncomment when UP-17662 is fixed
#resource uptime_contact test {
#  name     = "{{ petname 3 "-" }}"
#  sms_list = ["+44201234123"]
#}
// ---
# TODO: Uncomment when UP-17662 is fixed
#resource uptime_contact test {
#  name           = "{{ petname 3 "-" }}"
#  phonecall_list = ["+44201234123"]
#}
