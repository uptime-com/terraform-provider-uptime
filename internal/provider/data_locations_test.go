package provider

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLocationsDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_locations"),
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
				checkAnyLocationHasIP,
			),
		},
	}))
}

func TestAccLocationsDataSourceNewFields(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_locations"),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAnyLocationHasAddresses,
			),
		},
	}))
}

// The probe-servers list includes virtual routing locations such as the "Default Auto
// M2 Server" (location AUTO) that carry no fixed IPs, and one of them can sort to index
// 0. Assert that some location exposes addresses instead of pinning the checks to
// locations.0, which is otherwise flaky depending on ordering.
func checkAnyLocationHasIP(s *terraform.State) error {
	attrs, err := locationsAttrs(s)
	if err != nil {
		return err
	}
	count, _ := strconv.Atoi(attrs["locations.#"])
	for i := 0; i < count; i++ {
		if attrs[fmt.Sprintf("locations.%d.ip", i)] != "" {
			return nil
		}
	}
	return errors.New("no location exposes a non-empty ip")
}

func checkAnyLocationHasAddresses(s *terraform.State) error {
	attrs, err := locationsAttrs(s)
	if err != nil {
		return err
	}
	count, _ := strconv.Atoi(attrs["locations.#"])
	for i := 0; i < count; i++ {
		ipv4, _ := strconv.Atoi(attrs[fmt.Sprintf("locations.%d.ipv4_addresses.#", i)])
		ipv6, _ := strconv.Atoi(attrs[fmt.Sprintf("locations.%d.ipv6_addresses.#", i)])
		if ipv4 >= 1 && ipv6 >= 1 {
			return nil
		}
	}
	return errors.New("no location exposes both an IPv4 and an IPv6 address")
}

func locationsAttrs(s *terraform.State) (map[string]string, error) {
	rs, ok := s.RootModule().Resources["data.uptime_locations.test"]
	if !ok {
		return nil, errors.New("data.uptime_locations.test not found in state")
	}
	return rs.Primary.Attributes, nil
}
