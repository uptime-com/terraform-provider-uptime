package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestMaintenanceNotificationToAPIArgument(t *testing.T) {
	a := MaintenanceNotificationModelAdapter{}
	m := MaintenanceNotificationModel{
		ID:         types.Int64Value(101),
		ScheduleID: types.Int64Value(7),
		Offset:     types.Int64Value(-1800),
		Event:      types.StringValue("START"),
		ContactGroups: types.SetValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(10),
		}),
	}
	arg, err := a.ToAPIArgument(m)
	require.NoError(t, err)
	require.Equal(t, int64(7), arg.ScheduleID)
	require.Equal(t, int64(-1800), arg.Offset)
	require.Equal(t, "START", arg.Event)
	require.Equal(t, []int64{10}, arg.ContactGroups)
}

func TestMaintenanceNotificationToAPIArgumentEmptyContactGroups(t *testing.T) {
	a := MaintenanceNotificationModelAdapter{}
	m := MaintenanceNotificationModel{
		ScheduleID:    types.Int64Value(7),
		Offset:        types.Int64Value(0),
		Event:         types.StringValue("END"),
		ContactGroups: types.SetNull(types.Int64Type),
	}
	arg, err := a.ToAPIArgument(m)
	require.NoError(t, err)
	// contact_groups must serialize as [] not null.
	require.Equal(t, []int64{}, arg.ContactGroups)
}

func TestMaintenanceNotificationFromAPIResult(t *testing.T) {
	a := MaintenanceNotificationModelAdapter{}
	api := upapi.MaintenanceScheduleNotification{
		PK:            101,
		ScheduleID:    7,
		Offset:        -1800,
		Event:         "START",
		ContactGroups: []upapi.MaintenanceContactGroup{{PK: 10}},
		CreatedAt:     "2026-06-15T14:35:00Z",
		ModifiedAt:    "2026-06-15T14:35:00Z",
	}
	model, err := a.FromAPIResult(api)
	require.NoError(t, err)
	require.Equal(t, int64(101), model.ID.ValueInt64())
	require.Equal(t, int64(7), model.ScheduleID.ValueInt64())
	require.Equal(t, int64(-1800), model.Offset.ValueInt64())
	require.Equal(t, "START", model.Event.ValueString())
	require.Equal(t, []int64{10}, a.Slice(model.ContactGroups))
	require.False(t, model.CreatedAt.IsNull())
	require.False(t, model.ModifiedAt.IsNull())
}

func TestAccMaintenanceNotificationResource_Basic(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_maintenance_notification/basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_maintenance_notification.test", "offset", "-1800"),
				resource.TestCheckResourceAttr("uptime_maintenance_notification.test", "event", "START"),
				resource.TestCheckResourceAttrSet("uptime_maintenance_notification.test", "schedule_id"),
				resource.TestCheckResourceAttr("uptime_maintenance_notification.test", "contact_groups.#", "1"),
			),
		},
		{
			ResourceName:      "uptime_maintenance_notification.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}))
}
