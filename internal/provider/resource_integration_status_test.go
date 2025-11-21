package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationStatusResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	statuspageID := [2]string{
		"test-statuspage-id-12345", // 25 chars (min 24)
		"test-statuspage-id-23456", // 25 chars
	}
	apiID := [2]string{
		"test-api-id-123456789012345678901234", // 36 chars, meets min 36 requirement
		"test-api-id-234567890123456789012345", // 36 chars
	}
	apiKey := [2]string{
		"test-api-key-1234567890123456789012345678901234567890123456789012", // 65 chars (min 64)
		"test-api-key-2345678901234567890123456789012345678901234567890123", // 65 chars
	}
	metric := [2]string{
		"test-metric-123456789012", // 25 chars (min 24)
		"test-metric-234567890123", // 25 chars
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":          config.StringVariable(names[0]),
				"statuspage_id": config.StringVariable(statuspageID[0]),
				"api_id":        config.StringVariable(apiID[0]),
				"api_key":       config.StringVariable(apiKey[0]),
				"metric":        config.StringVariable(metric[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_status/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_status.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_status.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_status.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_status.test", "statuspage_id", statuspageID[0]),
				resource.TestCheckResourceAttr("uptime_integration_status.test", "api_id", apiID[0]),
				resource.TestCheckResourceAttr("uptime_integration_status.test", "api_key", apiKey[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":          config.StringVariable(names[1]),
				"statuspage_id": config.StringVariable(statuspageID[1]),
				"api_id":        config.StringVariable(apiID[1]),
				"api_key":       config.StringVariable(apiKey[1]),
				"metric":        config.StringVariable(metric[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_status/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_status.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_status.test", "statuspage_id", statuspageID[1]),
				resource.TestCheckResourceAttr("uptime_integration_status.test", "api_id", apiID[1]),
				resource.TestCheckResourceAttr("uptime_integration_status.test", "api_key", apiKey[1]),
			),
		},
	}))
}
