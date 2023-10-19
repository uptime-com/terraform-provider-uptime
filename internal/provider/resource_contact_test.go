package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactResource_EmailList(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/email_list"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
				"email_list": config.ListVariable(
					config.StringVariable("nobody@uptime.com"),
					config.StringVariable("noone@uptime.com"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_contact.test", "name", names[0]),

				resource.TestCheckResourceAttr("uptime_contact.test", "email_list.#", "2"),
				resource.TestCheckResourceAttr("uptime_contact.test", "email_list.0", "nobody@uptime.com"),
				resource.TestCheckResourceAttr("uptime_contact.test", "email_list.1", "noone@uptime.com"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/email_list"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
				"email_list": config.ListVariable(
					config.StringVariable("nobody@uptime.com"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_contact.test", "name", names[1]),

				resource.TestCheckResourceAttr("uptime_contact.test", "email_list.#", "1"),
				resource.TestCheckResourceAttr("uptime_contact.test", "email_list.0", "nobody@uptime.com"),
			),
		},
	}))
}

func TestAccContactResource_PhonecallList(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/phonecall_list"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
				"phonecall_list": config.ListVariable(
					config.StringVariable("+12065550100"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_contact.test", "name", names[0]),

				resource.TestCheckResourceAttr("uptime_contact.test", "phonecall_list.#", "1"),
				resource.TestCheckResourceAttr("uptime_contact.test", "phonecall_list.0", "+12065550100"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/phonecall_list"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
				"phonecall_list": config.ListVariable(
					config.StringVariable("+12065550100"),
					config.StringVariable("+441134960000"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_contact.test", "name", names[1]),

				resource.TestCheckResourceAttr("uptime_contact.test", "phonecall_list.#", "2"),
				resource.TestCheckResourceAttr("uptime_contact.test", "phonecall_list.0", "+12065550100"),
				resource.TestCheckResourceAttr("uptime_contact.test", "phonecall_list.1", "+441134960000"),
			),
		},
	}))
}

func TestAccContactResource_SMSList(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/sms_list"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
				"sms_list": config.ListVariable(
					config.StringVariable("+12065550100"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_contact.test", "name", names[0]),

				resource.TestCheckResourceAttr("uptime_contact.test", "sms_list.#", "1"),
				resource.TestCheckResourceAttr("uptime_contact.test", "sms_list.0", "+12065550100"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/sms_list"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
				"sms_list": config.ListVariable(
					config.StringVariable("+12065550100"),
					config.StringVariable("+441134960000"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_contact.test", "name", names[1]),

				resource.TestCheckResourceAttr("uptime_contact.test", "sms_list.#", "2"),
				resource.TestCheckResourceAttr("uptime_contact.test", "sms_list.0", "+12065550100"),
				resource.TestCheckResourceAttr("uptime_contact.test", "sms_list.1", "+441134960000"),
			),
		},
	}))
}
