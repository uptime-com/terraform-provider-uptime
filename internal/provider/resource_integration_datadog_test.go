package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationDatadogResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	apiKey := [2]string{
		"dd-api-key-000000000000000000000", // 32 chars
		"dd-api-key-111111111111111111111",
	}
	appKey := [2]string{
		"dd-app-key-00000000000000000000000000000", // 40 chars
		"dd-app-key-11111111111111111111111111111",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"api_key": config.StringVariable(apiKey[0]),
				"app_key": config.StringVariable(appKey[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_datadog/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_datadog.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_datadog.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_datadog.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_datadog.test", "api_key", apiKey[0]),
				resource.TestCheckResourceAttr("uptime_integration_datadog.test", "app_key", appKey[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"api_key": config.StringVariable(apiKey[1]),
				"app_key": config.StringVariable(appKey[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_datadog/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_datadog.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_datadog.test", "api_key", apiKey[1]),
				resource.TestCheckResourceAttr("uptime_integration_datadog.test", "app_key", appKey[1]),
			),
		},
	}))
}
