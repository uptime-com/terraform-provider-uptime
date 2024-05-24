package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckSSHResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"address": config.StringVariable("example.com"),
				"port":    config.IntegerVariable(80),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.com"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("example.net"),
				"port":    config.IntegerVariable(80),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.net"),
			),
		},
	}))
}

func TestAccCheckSSHResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckSSHResource_Interval(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"port":     config.IntegerVariable(80),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"port":     config.IntegerVariable(80),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckSSHResource_Locations(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"locations": config.ListVariable(
					config.StringVariable("US-CA-Los Angeles"),
					config.StringVariable("US-NY-New York"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.0", "US-CA-Los Angeles"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.1", "US-NY-New York"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"locations": config.ListVariable(
					config.StringVariable("Finland-Helsinki"),
					config.StringVariable("Switzerland-Zurich"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.0", "Finland-Helsinki"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.1", "Switzerland-Zurich"),
			),
		},
	}))
}

func TestAccCheckSSHResource_NumRetries(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"port":        config.IntegerVariable(80),
				"num_retries": config.IntegerVariable(3),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "num_retries", "3"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"port":        config.IntegerVariable(80),
				"num_retries": config.IntegerVariable(2),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "num_retries", "2"),
			),
		},
	}))
}

func TestAccCheckSSHResource_Sensitivity(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[0]),
				"address":     config.StringVariable("example.com"),
				"port":        config.IntegerVariable(80),
				"sensitivity": config.IntegerVariable(3),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/sensitivity"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.com"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "sensitivity", "3"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[1]),
				"address":     config.StringVariable("example.net"),
				"port":        config.IntegerVariable(80),
				"sensitivity": config.IntegerVariable(4),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/sensitivity"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.net"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "sensitivity", "4"),
			),
		},
	}))
}
