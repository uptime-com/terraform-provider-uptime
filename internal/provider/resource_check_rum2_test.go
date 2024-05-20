package provider

import (
	"sort"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckRUM2Resource(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(name),
				"address": config.StringVariable("example.com"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rum2/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "name", name),
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "address", "example.com"),
			),
		},
	}))
}

func TestAccCheckRUM2Resource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rum2/contact_groups"),
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(name),
				"address": config.StringVariable("example.com"),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "contact_groups.0", "nobody"),
			),
		},
	}))
}

func TestAccCheckRUM2Resource_UptimeSLA(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rum2/sla"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"address":    config.StringVariable("example.com"),
				"sla_uptime": config.StringVariable("0.88"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "sla_uptime", "0.88"),
			),
		},
	}))
}

func TestAccCheckRUM2Resource_Tags(t *testing.T) {
	name := petname.Generate(3, "-")
	tags := []string{
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
	}
	sort.Strings(tags)
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rum2/tags"),
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(name),
				"address": config.StringVariable("example.com"),
				"tags_create": config.SetVariable(
					config.StringVariable(tags[0]),
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
				"tags_use": config.SetVariable(
					config.StringVariable(tags[0]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "tags.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "tags.0", tags[0]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_rum2/tags"),
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(name),
				"address": config.StringVariable("example.com"),
				"tags_create": config.SetVariable(
					config.StringVariable(tags[0]),
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
				"tags_use": config.SetVariable(
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "tags.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "tags.0", tags[1]),
				resource.TestCheckResourceAttr("uptime_check_rum2.test", "tags.1", tags[2]),
			),
		},
	}))
}
