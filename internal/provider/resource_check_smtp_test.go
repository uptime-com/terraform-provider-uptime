package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckSMTPResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"address": config.StringVariable("example.com"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "address", "example.com"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("example.net"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "address", "example.net"),
			),
		},
	}))
}

func TestAccCheckSMTPResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckSMTPResource_Interval(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckSMTPResource_Locations(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("US-CA-Los Angeles"),
					config.StringVariable("US-NY-New York"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "locations.0", "US-CA-Los Angeles"),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "locations.1", "US-NY-New York"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("Israel-Tel Aviv"),
					config.StringVariable("Serbia-Belgrade"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "locations.0", "Israel-Tel Aviv"),
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "locations.1", "Serbia-Belgrade"),
			),
		},
	}))
}

func TestAccCheckSMTPResource_Port(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/port"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(5143),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "port", "5143"),
			),
		},
	}))
}

func TestAccCheckSMTPResource_NumRetries(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"num_retries": config.IntegerVariable(3),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "num_retries", "3"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_smtp/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"num_retries": config.IntegerVariable(2),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_smtp.test", "num_retries", "2"),
			),
		},
	}))
}
