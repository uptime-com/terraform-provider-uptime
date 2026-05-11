package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudStatusGroupsDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			Config: `data "uptime_cloudstatus_groups" "test" {}`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.uptime_cloudstatus_groups.test", "id", ""),
				resource.TestCheckResourceAttrSet("data.uptime_cloudstatus_groups.test", "groups.#"),
			),
		},
	}))
}

func TestAccCloudStatusGroupsDataSourceWithSearch(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			Config: `data "uptime_cloudstatus_groups" "test" { search = "amazon" }`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.uptime_cloudstatus_groups.test", "search", "amazon"),
				resource.TestCheckResourceAttrSet("data.uptime_cloudstatus_groups.test", "groups.#"),
			),
		},
	}))
}
