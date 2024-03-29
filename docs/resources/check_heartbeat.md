---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "uptime_check_heartbeat Resource - terraform-provider-uptime"
subcategory: ""
description: |-
  Monitor a periodic process, such as Cron, and issue alerts if the expected interval is exceeded
---

# uptime_check_heartbeat (Resource)

Monitor a periodic process, such as Cron, and issue alerts if the expected interval is exceeded



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `contact_groups` (Set of String)
- `include_in_global_metrics` (Boolean) Include this check in uptime/response time calculations for the dashboard and status pages
- `interval` (Number) The interval between checks in minutes
- `is_paused` (Boolean)
- `notes` (String)
- `sla` (Attributes) SLA related attributes (see [below for nested schema](#nestedatt--sla))
- `tags` (Set of String)

### Read-Only

- `heartbeat_url` (String) URL to send data to the check
- `id` (Number) The ID of this resource.
- `url` (String)

<a id="nestedatt--sla"></a>
### Nested Schema for `sla`

Optional:

- `latency` (String) The maximum average response time. Unit is mandatory (e.g. 1500ms or 1.5s or 1s500ms).
- `uptime` (String) The minimum uptime percentage. \nMust be a fraction with exactly 4 decimal places (e.g. 0.9995 for 99.95% uptime)


