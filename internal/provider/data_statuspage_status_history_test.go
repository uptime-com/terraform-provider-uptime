package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageStatusHistoryDataSource(t *testing.T) {
	spName := petname.Generate(3, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_statuspage_status_history"),
			ConfigVariables: config.Variables{
				"statuspage_name": config.StringVariable(spName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created status page
				resource.TestCheckResourceAttr("uptime_statuspage.test", "name", spName),
				// Check that status history data source exists and returns history (may be empty for new status page)
				resource.TestCheckOutput("history_count", "0"),
			),
		},
	}))
}
