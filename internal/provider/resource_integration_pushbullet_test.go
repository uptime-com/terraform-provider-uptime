package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationPushbulletResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	email := [2]string{
		"test1@example.com",
		"test2@example.com",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":  config.StringVariable(names[0]),
				"email": config.StringVariable(email[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_pushbullet/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_pushbullet.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_pushbullet.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_pushbullet.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_pushbullet.test", "email", email[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":  config.StringVariable(names[1]),
				"email": config.StringVariable(email[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_pushbullet/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_pushbullet.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_pushbullet.test", "email", email[1]),
			),
		},
	}))
}
