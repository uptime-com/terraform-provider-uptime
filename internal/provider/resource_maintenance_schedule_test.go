package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func mustRFC3339(t *testing.T, s string) timetypes.RFC3339 {
	t.Helper()
	v, diags := timetypes.NewRFC3339Value(s)
	require.False(t, diags.HasError(), "invalid RFC3339 %q: %v", s, diags)
	return v
}

func TestAccMaintenanceScheduleRRuleRequiresRrule(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{{
		Config: `resource "uptime_maintenance_schedule" "t" {
			name             = "x"
			schedule_type    = "RRULE"
			starts_at        = "2030-01-01T02:00:00Z"
			duration_minutes = 60
		}`,
		PlanOnly:    true,
		ExpectError: regexp.MustCompile(`schedule_type RRULE requires`),
	}}))
}

func TestAccMaintenanceScheduleOneOffRequiresDurationOrEnd(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{{
		Config: `resource "uptime_maintenance_schedule" "t" {
			name          = "x"
			schedule_type = "ONE_OFF"
			starts_at     = "2030-01-01T02:00:00Z"
		}`,
		PlanOnly:    true,
		ExpectError: regexp.MustCompile(`schedule_type ONE_OFF requires duration_minutes or ends_at`),
	}}))
}

func TestAccMaintenanceScheduleResource_RRule(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_maintenance_schedule/rrule"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "name", "tf-acc-rrule"),
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "schedule_type", "RRULE"),
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "rrule", "FREQ=WEEKLY;BYDAY=SA"),
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "duration_minutes", "120"),
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "is_active", "true"),
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "pause_checks_during_maintenance", "true"),
			),
		},
		{
			ResourceName:      "uptime_maintenance_schedule.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}))
}

func TestAccMaintenanceScheduleResource_OneOff(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_maintenance_schedule/one_off"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "name", "tf-acc-one-off"),
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "schedule_type", "ONE_OFF"),
				resource.TestCheckResourceAttr("uptime_maintenance_schedule.test", "duration_minutes", "120"),
				// Server computes ends_at from starts_at + duration.
				resource.TestCheckResourceAttrSet("uptime_maintenance_schedule.test", "ends_at"),
			),
		},
		{
			// No perpetual diff on the server-computed ends_at.
			ConfigDirectory: config.StaticDirectory("testdata/resource_maintenance_schedule/one_off"),
			PlanOnly:        true,
		},
		{
			ResourceName:      "uptime_maintenance_schedule.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}))
}

func TestMaintenanceScheduleToAPIArgument(t *testing.T) {
	a := MaintenanceScheduleModelAdapter{}
	m := MaintenanceScheduleModel{
		ID:                           types.Int64Value(7),
		Name:                         types.StringValue("x"),
		ScheduleType:                 types.StringValue("RRULE"),
		StartsAt:                     mustRFC3339(t, "2026-06-20T02:00:00Z"),
		EndsAt:                       timetypes.NewRFC3339Null(),
		RRule:                        types.StringValue("FREQ=WEEKLY;BYDAY=SA"),
		DurationMinutes:              types.Int64Value(120),
		IsActive:                     types.BoolValue(false),
		PauseChecksDuringMaintenance: types.BoolValue(true),
		Services:                     types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(42)}),
		Tags:                         types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(5)}),
	}
	arg, err := a.ToAPIArgument(m)
	require.NoError(t, err)
	require.Equal(t, "RRULE", arg.ScheduleType)
	require.Equal(t, int64(120), *arg.DurationMinutes)
	require.False(t, arg.IsActive)
	require.True(t, arg.PauseChecksDuringMaintenance)
	require.Equal(t, []int64{42}, arg.Services)
	require.Equal(t, []int64{5}, arg.Tags)
	require.Nil(t, arg.EndsAt)
}

func TestMaintenanceScheduleToAPIArgumentOneOff(t *testing.T) {
	a := MaintenanceScheduleModelAdapter{}
	m := MaintenanceScheduleModel{
		Name:                         types.StringValue("x"),
		ScheduleType:                 types.StringValue("ONE_OFF"),
		StartsAt:                     mustRFC3339(t, "2026-06-20T02:00:00Z"),
		EndsAt:                       mustRFC3339(t, "2026-06-20T04:00:00Z"),
		RRule:                        types.StringNull(),
		DurationMinutes:              types.Int64Null(),
		IsActive:                     types.BoolValue(true),
		PauseChecksDuringMaintenance: types.BoolValue(false),
		Services:                     types.SetNull(types.Int64Type),
		Tags:                         types.SetNull(types.Int64Type),
	}
	arg, err := a.ToAPIArgument(m)
	require.NoError(t, err)
	require.Equal(t, "ONE_OFF", arg.ScheduleType)
	require.NotNil(t, arg.EndsAt)
	require.Equal(t, "2026-06-20T04:00:00Z", *arg.EndsAt)
	require.Nil(t, arg.DurationMinutes)
	// Empty targets must serialize as [] not null.
	require.Equal(t, []int64{}, arg.Services)
	require.Equal(t, []int64{}, arg.Tags)
}

func TestMaintenanceScheduleFromAPIResult(t *testing.T) {
	a := MaintenanceScheduleModelAdapter{}
	dur := int64(120)
	api := upapi.MaintenanceSchedule{
		PK:                           7,
		Name:                         "x",
		ScheduleType:                 "RRULE",
		StartsAt:                     "2026-06-20T02:00:00Z",
		EndsAt:                       "2026-09-20T02:00:00Z",
		RRule:                        "FREQ=WEEKLY;BYDAY=SA",
		DurationMinutes:              &dur,
		IsActive:                     true,
		PauseChecksDuringMaintenance: true,
		Services:                     []upapi.MaintenanceService{{PK: 42}},
		Tags:                         []upapi.MaintenanceTag{{PK: 5}},
		CreatedAt:                    "2026-06-15T14:32:00Z",
		ModifiedAt:                   "2026-06-15T14:32:00Z",
	}
	model, err := a.FromAPIResult(api)
	require.NoError(t, err)
	require.Equal(t, int64(7), model.ID.ValueInt64())
	require.Equal(t, "FREQ=WEEKLY;BYDAY=SA", model.RRule.ValueString())
	require.Equal(t, int64(120), model.DurationMinutes.ValueInt64())
	require.False(t, model.StartsAt.IsNull())
	require.False(t, model.EndsAt.IsNull())

	serviceIDs := a.Slice(model.Services)
	require.Equal(t, []int64{42}, serviceIDs)
	tagIDs := a.Slice(model.Tags)
	require.Equal(t, []int64{5}, tagIDs)
}

func TestMaintenanceScheduleFromAPIResultOneOff(t *testing.T) {
	a := MaintenanceScheduleModelAdapter{}
	api := upapi.MaintenanceSchedule{
		PK:           8,
		Name:         "y",
		ScheduleType: "ONE_OFF",
		StartsAt:     "2026-06-20T02:00:00Z",
		EndsAt:       "2026-06-20T04:00:00Z",
		RRule:        "", // ONE_OFF has no rrule
		IsActive:     true,
	}
	model, err := a.FromAPIResult(api)
	require.NoError(t, err)
	// empty rrule -> Null (avoids drift)
	require.True(t, model.RRule.IsNull())
	// nil duration -> Null
	require.True(t, model.DurationMinutes.IsNull())
	require.False(t, model.EndsAt.IsNull())
}
