package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckGroupResource(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_group/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_group.test", "name", name),
			),
		},
	}))
}

func TestAccCheckGroupResource_Config(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"check_name": config.StringVariable(name),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_group/config"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_group.test", "name", name),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.services.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.services.0", name),
			),
		},
	}))
}

func TestAccCheckGroupResource_ResponseTime(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_group/response_time"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_group.test", "name", name),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.response_time.check_type", "HTTP"),
			),
		},
	}))
}

func TestAccCheckGroupResource_PercentCalculation(t *testing.T) {
	names := []string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":                       config.StringVariable(names[0]),
				"uptime_percent_calculation": config.StringVariable("UP_DOWN_STATES"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_group/percent_calculation"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_group.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.response_time.check_type", "HTTP"),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.uptime_percent_calculation", "UP_DOWN_STATES"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":                       config.StringVariable(names[1]),
				"uptime_percent_calculation": config.StringVariable("AVERAGE"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_group/percent_calculation"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_group.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.response_time.check_type", "HTTP"),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.uptime_percent_calculation", "AVERAGE"),
			),
		},
	}))
}

func TestAccCheckGroupResource_DownCondition(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":           config.StringVariable(name),
				"down_condition": config.StringVariable("TEN_PCT"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_group/down_condition"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_group.test", "name", name),
				resource.TestCheckResourceAttr("uptime_check_group.test", "config.down_condition", "TEN_PCT"),
			),
		},
	}))
}
