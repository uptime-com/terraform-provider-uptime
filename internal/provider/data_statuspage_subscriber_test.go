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

func TestAccStatusPageSubscriberDataSource(t *testing.T) {
	spName := petname.Generate(3, "-")
	email := fmt.Sprintf("%s@example.com", petname.Generate(2, "-"))

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_statuspage_subscriber"),
			ConfigVariables: config.Variables{
				"statuspage_name": config.StringVariable(spName),
				"email":           config.StringVariable(email),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created subscriber
				resource.TestCheckResourceAttr("uptime_statuspage_subscriber.test", "target", email),
				// Check that filtered output contains exactly 1 subscriber
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
						return fmt.Errorf("expected exactly 1 filtered subscriber, got %d", count)
					}
					return nil
				},
				// Check the filtered subscriber has the correct email
				resource.TestCheckOutput("filtered_subscriber_target", email),
			),
		},
	}))
}
