package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckMaintenanceResource_Basic(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_maintenance/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "state", "SUPPRESSED"),
				resource.TestCheckResourceAttr("uptime_check_maintenance.test", "schedule.#", "0"),
			),
		},
	}))
}
