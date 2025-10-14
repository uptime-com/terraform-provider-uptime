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

func TestAccStatusPageIncidentDataSource(t *testing.T) {
	spName := petname.Generate(3, "-")
	incName := petname.Generate(3, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_statuspage_incident"),
			ConfigVariables: config.Variables{
				"statuspage_name": config.StringVariable(spName),
				"incident_name":   config.StringVariable(incName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created incident
				resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "name", incName),
				// Check that filtered output contains exactly 1 incident
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
						return fmt.Errorf("expected exactly 1 filtered incident, got %d", count)
					}
					return nil
				},
				// Check the filtered incident has the correct name
				resource.TestCheckOutput("filtered_incident_name", incName),
			),
		},
	}))
}
