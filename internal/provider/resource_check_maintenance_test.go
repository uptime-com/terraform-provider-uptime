package provider

import (
	"regexp"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckMaintenanceResource_Basic(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":  config.StringVariable(name),
				"state": config.StringVariable("SUPPRESSED"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "state", "SUPPRESSED"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "0"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":  config.StringVariable(name),
				"state": config.StringVariable("ACTIVE"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "state", "ACTIVE"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "0"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/_basic"),
			ConfigVariables: config.Variables{
				"name":  config.StringVariable(name),
				"state": config.StringVariable("__DOES_NOT_EXIST__"),
			},
			ExpectError: regexp.MustCompile("value must be one of"),
		},
	}))
}

func TestAccCheckMaintenanceResource_Weekly(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":      config.StringVariable(name),
				"from_time": config.StringVariable("10:00:00"),
				"to_time":   config.StringVariable("11:00:00"),
				"weekdays":  config.ListVariable(config.IntegerVariable(0), config.IntegerVariable(6)),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/weekly"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.weekdays.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.weekdays.0", "0"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.weekdays.1", "6"),
			),
		},
	}))
}

func TestAccCheckMaintenanceResource_Monthly(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":      config.StringVariable(name),
				"from_time": config.StringVariable("13:01:03"),
				"to_time":   config.StringVariable("14:59:59"),
				"monthday":  config.IntegerVariable(3),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/monthly"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.monthday", "3"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":          config.StringVariable(name),
				"from_time":     config.StringVariable("15:31:58"),
				"to_time":       config.StringVariable("23:50:09"),
				"monthday":      config.IntegerVariable(0),
				"monthday_from": config.IntegerVariable(19),
				"monthday_to":   config.IntegerVariable(23),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/monthly"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.monthday_from", "19"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.monthday_to", "23"),
			),
		},
	}))
}

func TestAccCheckMaintenanceResource_Once(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"once_start_date": config.StringVariable("2024-10-23T13:01:03Z"),
				"once_end_date":   config.StringVariable("2024-10-25T14:59:59Z"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/once"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.once_start_date", "2024-10-23T13:01:03Z"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.once_end_date", "2024-10-25T14:59:59Z"),
			),
		},
	}))
}

func TestAccCheckMaintenanceResource_All(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":              config.StringVariable(name),
				"from_time":         config.StringVariable("13:01:03"),
				"to_time":           config.StringVariable("14:59:59"),
				"weekdays":          config.ListVariable(config.IntegerVariable(0), config.IntegerVariable(6)),
				"monthly_from_time": config.StringVariable("15:01:03"),
				"monthly_to_time":   config.StringVariable("16:59:59"),
				"monthday":          config.IntegerVariable(0),
				"monthday_from":     config.IntegerVariable(19),
				"monthday_to":       config.IntegerVariable(23),
				"once_start_date":   config.StringVariable("2024-10-23T13:01:03Z"),
				"once_end_date":     config.StringVariable("2024-10-25T14:59:59Z"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/all"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "3"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.type", "WEEKLY"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.from_time", "13:01:03"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.to_time", "14:59:59"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.weekdays.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.weekdays.0", "0"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.0.weekdays.1", "6"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.1.type", "MONTHLY"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.1.from_time", "15:01:03"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.1.to_time", "16:59:59"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.1.monthday_from", "19"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.1.monthday_to", "23"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.2.type", "ONCE"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.2.once_start_date", "2024-10-23T13:01:03Z"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.2.once_end_date", "2024-10-25T14:59:59Z"),
			),
		},
	}))
}
