package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckAPIResource(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_api.test", "locations.#", "2"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "name", names[1]),
			),
		},
	}))
}

func TestAccCheckAPIResource_ContactGroups(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_api.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_api.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_api.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckAPIResource_Interval(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckAPIResource_Locations(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("US East"),
					config.StringVariable("US West"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_api.test", "locations.0", "US East"),
				resource.TestCheckResourceAttr("uptime_check_api.test", "locations.1", "US West"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("Austria"),
					config.StringVariable("Serbia"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_api.test", "locations.0", "Austria"),
				resource.TestCheckResourceAttr("uptime_check_api.test", "locations.1", "Serbia"),
			),
		},
	}))
}

func TestAccCheckAPIResource_ResponseTimeSLA(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/response_time_sla"),
			ConfigVariables: config.Variables{
				"name":              config.StringVariable(name),
				"response_time_sla": config.StringVariable("100ms"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "response_time_sla", "100ms"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_api/response_time_sla"),
			ConfigVariables: config.Variables{
				"name":              config.StringVariable(name),
				"response_time_sla": config.StringVariable("200ms"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_api.test", "response_time_sla", "200ms"),
			),
		},
	}))
}
