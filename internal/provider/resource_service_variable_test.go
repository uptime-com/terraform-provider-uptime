package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceVariableResource(t *testing.T) {
	credentialName := petname.Generate(3, "-")
	password := petname.Generate(1, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_basic"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
				"variable_name":   config.StringVariable("api_password"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "api_password"),
				resource.TestCheckResourceAttr("uptime_service_variable.test", "property_name", "password"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "service_id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "credential_id"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_basic"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
				"variable_name":   config.StringVariable("api_key"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "api_key"),
				resource.TestCheckResourceAttr("uptime_service_variable.test", "property_name", "password"),
			),
		},
	}))
}

func TestAccServiceVariableResource_WithDataSource(t *testing.T) {
	credentialName := petname.Generate(3, "-")
	password := petname.Generate(1, "-")
	variableName := "api_password"

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/datasource_step1"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", credentialName),
				resource.TestCheckResourceAttr("uptime_credential.test", "credential_type", "BASIC"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/datasource_step2"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
				"variable_name":   config.StringVariable(variableName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check credential exists
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", credentialName),
				// Check datasource works
				resource.TestCheckResourceAttrSet("data.uptime_credentials.all", "credentials.#"),
				// Check service variable is created
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", variableName),
				resource.TestCheckResourceAttr("uptime_service_variable.test", "property_name", "password"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "service_id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "credential_id"),
				// Verify the credential_id matches the one from datasource
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.test", "id",
				),
			),
		},
	}))
}
