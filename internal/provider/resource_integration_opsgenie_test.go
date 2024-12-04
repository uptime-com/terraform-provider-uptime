package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationOpsgenieResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	apiEndpoint := [2]string{
		"https://api.opsgenie.com/v1/json/uptime1",
		"https://api.opsgenie.com/v1/json/uptime2",
	}
	apiKey := [2]string{
		"16c8bfe0-b219-11ef-a5da-4b9e62fe7439",
		"274d1848-b219-11ef-a468-9f18e59bdc97",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[0]),
				"api_endpoint": config.StringVariable(apiEndpoint[0]),
				"api_key":      config.StringVariable(apiKey[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_opsgenie/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_opsgenie.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_opsgenie.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_opsgenie.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_opsgenie.test", "api_endpoint", apiEndpoint[0]),
				resource.TestCheckResourceAttr("uptime_integration_opsgenie.test", "api_key", apiKey[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[1]),
				"api_endpoint": config.StringVariable(apiEndpoint[1]),
				"api_key":      config.StringVariable(apiKey[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_opsgenie/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_opsgenie.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_opsgenie.test", "api_endpoint", apiEndpoint[1]),
				resource.TestCheckResourceAttr("uptime_integration_opsgenie.test", "api_key", apiKey[1]),
			),
		},
	}))
}
