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

func TestAccStatusPageMetricDataSource(t *testing.T) {
	spName := petname.Generate(3, "-")
	metricName := petname.Generate(3, "-")
	checkName := petname.Generate(3, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_statuspage_metric"),
			ConfigVariables: config.Variables{
				"statuspage_name": config.StringVariable(spName),
				"metric_name":     config.StringVariable(metricName),
				"check_name":      config.StringVariable(checkName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created metric
				resource.TestCheckResourceAttr("uptime_statuspage_metric.test", "name", metricName),
				// Check that filtered output contains exactly 1 metric
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
						return fmt.Errorf("expected exactly 1 filtered metric, got %d", count)
					}
					return nil
				},
				// Check the filtered metric has the correct name
				resource.TestCheckOutput("filtered_metric_name", metricName),
			),
		},
	}))
}
