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
