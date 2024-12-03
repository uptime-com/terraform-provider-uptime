---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "uptime_statuspage_incident Resource - terraform-provider-uptime"
subcategory: ""
description: |-
  Status page incident or maintenance window resource
---

# uptime_statuspage_incident (Resource)

Status page incident or maintenance window resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `statuspage_id` (Number)
- `updates` (Attributes Set) (see [below for nested schema](#nestedatt--updates))

### Optional

- `affected_components` (Attributes Set) (see [below for nested schema](#nestedatt--affected_components))
- `ends_at` (String)
- `incident_type` (String)
- `include_in_global_metrics` (Boolean)
- `notify_subscribers` (Boolean)
- `send_maintenance_start_notification` (Boolean)
- `starts_at` (String)
- `update_component_status` (Boolean)

### Read-Only

- `id` (Number) The ID of this resource.
- `url` (String)

<a id="nestedatt--updates"></a>
### Nested Schema for `updates`

Optional:

- `description` (String)
- `incident_state` (String)

Read-Only:

- `id` (Number)


<a id="nestedatt--affected_components"></a>
### Nested Schema for `affected_components`

Required:

- `component_id` (Number)

Optional:

- `status` (String)

Read-Only:

- `id` (Number)

