package provider

import (
	"regexp"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialResource(t *testing.T) {
	names := [3]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	passwords := [3]string{
		petname.Generate(1, "-"),
		petname.Generate(1, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_credential/_basic"),
			ConfigVariables: config.Variables{
				"display_name":    config.StringVariable(names[0]),
				"credential_type": config.StringVariable("BASIC"),
				"password":        config.StringVariable(passwords[0]),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", names[0]),
				resource.TestCheckResourceAttr("uptime_credential.test", "credential_type", "BASIC"),
				resource.TestCheckResourceAttr("uptime_credential.test", "secret.password", passwords[0]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_credential/_basic"),
			ConfigVariables: config.Variables{
				"display_name":    config.StringVariable(names[1]),
				"credential_type": config.StringVariable("BASIC"),
				"password":        config.StringVariable(passwords[1]),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", names[1]),
				resource.TestCheckResourceAttr("uptime_credential.test", "credential_type", "BASIC"),
				resource.TestCheckResourceAttr("uptime_credential.test", "secret.password", passwords[1]),
			),
		},
	}))
}

func TestAccCredentialResource_Validation(t *testing.T) {
	names := [3]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	passwords := [3]string{
		petname.Generate(1, "-"),
		petname.Generate(1, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_credential/validation"),
			ConfigVariables: config.Variables{
				"display_name":    config.StringVariable(names[0]),
				"credential_type": config.StringVariable("BASIC"),
				"password":        config.StringVariable(passwords[0]),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", names[0]),
				resource.TestCheckResourceAttr("uptime_credential.test", "credential_type", "BASIC"),
				resource.TestCheckResourceAttr("uptime_credential.test", "secret.password", passwords[0]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_credential/validation"),
			ConfigVariables: config.Variables{
				"display_name":    config.StringVariable(names[2]),
				"credential_type": config.StringVariable("TOKEN"),
				"token":           config.StringVariable(passwords[2]),
			},
			ExpectError: regexp.MustCompile("When credential_type is TOKEN, only the secret field should be set."),
		},
	}))
}
