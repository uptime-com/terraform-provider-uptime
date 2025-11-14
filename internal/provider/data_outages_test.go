package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOutagesDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			Config: `data "uptime_outages" "test" {}`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.uptime_outages.test", "id", ""),
			),
		},
	}))
}
