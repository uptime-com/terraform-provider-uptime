package provider

import (
	"fmt"
	"strconv"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccStatusPageComponentDataSource(t *testing.T) {
	spName := petname.Generate(3, "-")
	compName := petname.Generate(3, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_statuspage_component"),
			ConfigVariables: config.Variables{
				"statuspage_name": config.StringVariable(spName),
				"component_name":  config.StringVariable(compName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created component
				resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", compName),
				// Check that filtered output contains exactly 1 component
				func(s *terraform.State) error {
					output, ok := s.RootModule().Outputs["filtered_count"]
					if !ok {
						return fmt.Errorf("filtered_count output not found")
					}
					count, err := strconv.Atoi(output.Value.(string))
					if err != nil {
						return fmt.Errorf("failed to parse filtered_count: %w", err)
					}
					if count != 1 {
						return fmt.Errorf("expected exactly 1 filtered component, got %d", count)
					}
					return nil
				},
				// Check the filtered component has the correct name
				resource.TestCheckOutput("filtered_component_name", compName),
			),
		},
	}))
}
