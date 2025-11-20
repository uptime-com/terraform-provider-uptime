package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageCurrentStatusDataSource(t *testing.T) {
	spName := petname.Generate(3, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_statuspage_current_status"),
			ConfigVariables: config.Variables{
				"statuspage_name": config.StringVariable(spName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created status page
				resource.TestCheckResourceAttr("uptime_statuspage.test", "name", spName),
				// Check that current status data source returns data
				resource.TestCheckResourceAttr("data.uptime_statuspage_current_status.test", "global_is_operational", "true"),
				resource.TestCheckOutput("components_count", "0"),
			),
		},
	}))
}
