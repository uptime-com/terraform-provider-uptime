package provider

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testRenderSnippet(t, "data_locations.tf", 0, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.uptime_locations.test", "locations.#",
						func(value string) error {
							n, err := strconv.Atoi(value)
							if err != nil {
								return fmt.Errorf("failed to parse locations count: %w", err)
							}
							if n < 3 {
								return errors.New("expected at least 3 locations")
							}
							return nil
						},
					),
				),
			},
		},
	})
}
