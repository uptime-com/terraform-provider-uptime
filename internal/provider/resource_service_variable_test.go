package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestServiceVariablePreservePlanValues verifies that Required attributes the update
// endpoint may omit (credential_id, property_name, variable_name) are restored from the
// plan when the API result comes back empty, preventing "Provider produced inconsistent
// result after apply".
func TestServiceVariablePreservePlanValues(t *testing.T) {
	adapter := ServiceVariableResourceModelAdapter{}

	plan := &ServiceVariableResourceModel{
		CredentialID: types.Int64Value(28966),
		PropertyName: types.StringValue("secret"),
		VariableName: types.StringValue("token_raven_token"),
	}

	t.Run("restores omitted values", func(t *testing.T) {
		result := &ServiceVariableResourceModel{
			CredentialID: types.Int64Value(0),
			PropertyName: types.StringValue(""),
			VariableName: types.StringValue(""),
		}
		got := adapter.PreservePlanValues(result, plan)
		if got.CredentialID.ValueInt64() != 28966 {
			t.Errorf("credential_id: got %d, want 28966", got.CredentialID.ValueInt64())
		}
		if got.PropertyName.ValueString() != "secret" {
			t.Errorf("property_name: got %q, want \"secret\"", got.PropertyName.ValueString())
		}
		if got.VariableName.ValueString() != "token_raven_token" {
			t.Errorf("variable_name: got %q, want \"token_raven_token\"", got.VariableName.ValueString())
		}
	})

	t.Run("keeps non-empty result values", func(t *testing.T) {
		result := &ServiceVariableResourceModel{
			CredentialID: types.Int64Value(28965),
			PropertyName: types.StringValue("password"),
			VariableName: types.StringValue("other_name"),
		}
		got := adapter.PreservePlanValues(result, plan)
		if got.CredentialID.ValueInt64() != 28965 {
			t.Errorf("credential_id: got %d, want 28965", got.CredentialID.ValueInt64())
		}
		if got.PropertyName.ValueString() != "password" {
			t.Errorf("property_name: got %q, want \"password\"", got.PropertyName.ValueString())
		}
		if got.VariableName.ValueString() != "other_name" {
			t.Errorf("variable_name: got %q, want \"other_name\"", got.VariableName.ValueString())
		}
	})
}

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
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.test", "id",
				),
			),
		},
	}))
}

// TestAccServiceVariableResource_TokenSwap reproduces the customer scenario: a TOKEN
// credential service variable whose credential_id changes while variable_name stays
// constant.
func TestAccServiceVariableResource_TokenSwap(t *testing.T) {
	credentialNameA := petname.Generate(3, "-")
	credentialNameB := petname.Generate(3, "-")
	token := petname.Generate(2, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_token_swap"),
			ConfigVariables: config.Variables{
				"credential_name_a": config.StringVariable(credentialNameA),
				"credential_name_b": config.StringVariable(credentialNameB),
				"token":             config.StringVariable(token),
				"use_b":             config.BoolVariable(false),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "token_raven_token"),
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.a", "id",
				),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_token_swap"),
			ConfigVariables: config.Variables{
				"credential_name_a": config.StringVariable(credentialNameA),
				"credential_name_b": config.StringVariable(credentialNameB),
				"token":             config.StringVariable(token),
				"use_b":             config.BoolVariable(true),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "token_raven_token"),
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.b", "id",
				),
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
