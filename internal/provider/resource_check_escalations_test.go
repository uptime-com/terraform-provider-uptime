package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccCheckEscalationsResource(t *testing.T) {
	names := [3]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_escalations/_basic"),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_check_escalations.test", "check_id"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.0.wait_time", "300"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.0.num_repeats", "3"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.0.contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.1.wait_time", "600"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.1.num_repeats", "0"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.1.contact_groups.#", "1"),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPreRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPostRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
			},
		},
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_escalations/_update"),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_check_escalations.test", "check_id"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.0.wait_time", "180"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.0.num_repeats", "5"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.0.contact_groups.#", "1"),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPreRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPostRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
			},
		},
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[2]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_escalations/_empty"),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_check_escalations.test", "check_id"),
				resource.TestCheckResourceAttr("uptime_check_escalations.test", "escalations.#", "0"),
			),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPreRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPostRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
			},
		},
	}))
}
