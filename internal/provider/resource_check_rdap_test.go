package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckRDAPResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	expectString := `
	Again in the country! A life full of pleasure:  
	I shoot; I write verses in solitude deep;  
	And yesterday, searching the moorland for treasure,  
	I came to a cowshed, turned in, fell asleep.
	`
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rdap/_basic"),
			ConfigVariables: config.Variables{
				"name":           config.StringVariable(names[0]),
				"contact_groups": config.ListVariable(config.StringVariable("nobody")),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "name", names[0]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rdap/_basic"),
			ConfigVariables: config.Variables{
				"name":           config.StringVariable(names[1]),
				"address":        config.StringVariable("example.net"),
				"expect_string":  config.StringVariable(expectString),
				"contact_groups": config.ListVariable(config.StringVariable("nobody")),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "address", "example.net"),
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "expect_string", expectString),
			),
		},
	}))
}

func TestAccCheckRDAPResource_Threshold(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rdap/threshold"),
			ConfigVariables: config.Variables{
				"name":           config.StringVariable(name),
				"threshold":      config.IntegerVariable(1),
				"contact_groups": config.ListVariable(config.StringVariable("nobody")),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "threshold", "1"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rdap/threshold"),
			ConfigVariables: config.Variables{
				"name":           config.StringVariable(name),
				"threshold":      config.IntegerVariable(2),
				"contact_groups": config.ListVariable(config.StringVariable("nobody")),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rdap.test", "threshold", "2"),
			),
		},
	}))
}
