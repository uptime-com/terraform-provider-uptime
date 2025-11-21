package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationPushoverResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	user := [2]string{
		"test-user-key-1",
		"test-user-key-2",
	}
	priority := [2]int64{
		1, // High priority
		2, // Emergency priority
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(names[0]),
				"user":     config.StringVariable(user[0]),
				"priority": config.IntegerVariable(priority[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_pushover/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_pushover.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_pushover.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_pushover.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_pushover.test", "user", user[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(names[1]),
				"user":     config.StringVariable(user[1]),
				"priority": config.IntegerVariable(priority[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_pushover/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_pushover.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_pushover.test", "user", user[1]),
			),
		},
	}))
}
