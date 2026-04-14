package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckCloudStatusResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[0]),
				"service_name": config.StringVariable("aws-ec2-us-east-1"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "service_name", "aws-ec2-us-east-1"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[1]),
				"service_name": config.StringVariable("aws-ec2-us-east-2"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "service_name", "aws-ec2-us-east-2"),
			),
		},
	}))
}

func TestAccCheckCloudStatusResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckCloudStatusResource_Group(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/group"),
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"group":           config.IntegerVariable(1),
				"monitoring_type": config.StringVariable("ALL"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "monitoring_type", "ALL"),
			),
		},
	}))
}
