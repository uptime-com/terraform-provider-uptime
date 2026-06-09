package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactResource_EmailList(t *testing.T) {
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
			ResourceName:      "uptime_contact.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}))
}

func TestAccContactResource_PhonecallList(t *testing.T) {
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

// TestAccContactResource_IntegrationNoDrift reproduces SYS-1264: when an
// integration references a contact group, the backend links the integration
// back into the contact's `integrations` field. The contact must not churn
// that server-managed field on subsequent plans.
func TestAccContactResource_IntegrationNoDrift(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/integration_drift"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_contact.test", "name", name),
				resource.TestCheckResourceAttr("uptime_integration_slack.test", "contact_groups.0", name),
			),
		},
		{
			// Same config: must produce an empty plan. With the old empty-set
			// Default on `integrations`, the server-added association churned
			// here as `~ integrations = [- "<name>"]`.
			ConfigDirectory: config.StaticDirectory("testdata/resource_contact/integration_drift"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
			},
			PlanOnly: true,
		},
	}))
}

func TestAccContactResource_SMSList(t *testing.T) {
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
