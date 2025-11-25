package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationKlipfolioResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	apiKey := [2]string{
		"test-api-key-123456789012345678901234567", // 40 chars, meets max 40 requirement
		"test-api-key-234567890123456789012345678", // 40 chars
	}
	dataSourceName := [2]string{
		"test-datasource-1",
		"test-datasource-2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":             config.StringVariable(names[0]),
				"api_key":          config.StringVariable(apiKey[0]),
				"data_source_name": config.StringVariable(dataSourceName[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_klipfolio/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_klipfolio.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_klipfolio.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_klipfolio.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_klipfolio.test", "api_key", apiKey[0]),
				resource.TestCheckResourceAttr("uptime_integration_klipfolio.test", "data_source_name", dataSourceName[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":             config.StringVariable(names[1]),
				"api_key":          config.StringVariable(apiKey[1]),
				"data_source_name": config.StringVariable(dataSourceName[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_klipfolio/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_klipfolio.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_klipfolio.test", "api_key", apiKey[1]),
				resource.TestCheckResourceAttr("uptime_integration_klipfolio.test", "data_source_name", dataSourceName[1]),
			),
		},
	}))
}
