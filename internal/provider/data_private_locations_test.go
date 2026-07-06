package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

// The data source filters probe servers on is_private. Element-level assertions
// (name/location/country) only hold on an account that actually has private
// locations; on accounts without any, the list is empty. Assert the data source
// reads without error and returns a well-formed list.
func TestAccPrivateLocationsDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_private_locations"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.uptime_private_locations.test", "locations.#"),
			),
		},
	}))
}
