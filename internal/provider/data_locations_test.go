package provider

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocationsDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_locations"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrWith("data.uptime_locations.test", "locations.0.ip",
					func(value string) error {
						if value == "" {
							return errors.New("host_ip property is empty")
						}
						return nil
					},
				),
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
	}))
}

func TestAccLocationsDataSourceNewFields(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_locations"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrWith("data.uptime_locations.test", "locations.0.ipv4_addresses.#",
					func(value string) error {
						n, err := strconv.Atoi(value)
						if err != nil {
							return fmt.Errorf("failed to parse ipv4_addresses count: %w", err)
						}
						if n < 1 {
							return errors.New("expected at least 1 IPv4 address")
						}
						return nil
					},
				),
				resource.TestCheckResourceAttrSet("data.uptime_locations.test", "locations.0.ipv4_addresses.0"),
				resource.TestCheckResourceAttrWith("data.uptime_locations.test", "locations.0.ipv6_addresses.#",
					func(value string) error {
						n, err := strconv.Atoi(value)
						if err != nil {
							return fmt.Errorf("failed to parse ipv6_addresses count: %w", err)
						}
						if n < 1 {
							return errors.New("expected at least 1 IPv6 address")
						}
						return nil
					},
				),
			),
		},
	}))
}
