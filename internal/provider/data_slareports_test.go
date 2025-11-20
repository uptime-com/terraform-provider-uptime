package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSLAReportsDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			Config: `data "uptime_sla_reports" "test" {}`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.uptime_sla_reports.test", "id", ""),
			),
		},
	}))
}
