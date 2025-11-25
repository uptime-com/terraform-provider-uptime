package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationMicrosoftTeamsResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	webhookURL := [2]string{
		"https://outlook.office.com/webhook/test1",
		"https://outlook.office.com/webhook/test2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[0]),
				"webhook_url": config.StringVariable(webhookURL[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_microsoft_teams/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_microsoft_teams.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_microsoft_teams.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_microsoft_teams.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_microsoft_teams.test", "webhook_url", webhookURL[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[1]),
				"webhook_url": config.StringVariable(webhookURL[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_microsoft_teams/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_microsoft_teams.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_microsoft_teams.test", "webhook_url", webhookURL[1]),
			),
		},
	}))
}
