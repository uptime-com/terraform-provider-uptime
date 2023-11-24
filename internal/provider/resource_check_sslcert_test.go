package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckSSLCertResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"address": config.StringVariable("example.com"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_sslcert/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "address", "example.com"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("example.net"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_sslcert/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "address", "example.net"),
			),
		},
	}))
}

func TestAccCheckSSLCertResource_Config_CRL(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_sslcert/config/crl"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"config_crl": config.BoolVariable(false),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "config.crl", "false"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_sslcert/config/crl"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"config_crl": config.BoolVariable(true),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "config.crl", "true"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_sslcert/config/crl"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"config_crl": config.BoolVariable(false),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "config.crl", "false"),
			),
		},
	}))
}

func TestAccCheckSSLCertResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_sslcert/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_sslcert/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_sslcert.test", "contact_groups.1", "noone"),
			),
		},
	}))
}
