package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudStatusServicesDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			Config: `
data "uptime_cloudstatus_groups" "groups" {
  search = "amazon"
}

data "uptime_cloudstatus_services" "services" {
  group = try(tostring(data.uptime_cloudstatus_groups.groups.groups[0].id), "1")
}
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.uptime_cloudstatus_services.services", "services.#"),
			),
		},
	}))
}

func TestAccCloudStatusServicesDataSourceFilterByGroupName(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			Config: `data "uptime_cloudstatus_services" "test" { group = "Amazon" }`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.uptime_cloudstatus_services.test", "group", "Amazon"),
				resource.TestCheckResourceAttrSet("data.uptime_cloudstatus_services.test", "services.#"),
			),
		},
	}))
}
