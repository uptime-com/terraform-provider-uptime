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

func TestAccStatusPageUserDataSource(t *testing.T) {
	t.Skip("Skipping test as API is broken: GET/DELETE methods")
	spName := petname.Generate(3, "-")
	email := fmt.Sprintf("%s@example.com", petname.Generate(2, "-"))
	firstName := petname.Generate(1, "-")
	lastName := petname.Generate(1, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_statuspage_user"),
			ConfigVariables: config.Variables{
				"statuspage_name": config.StringVariable(spName),
				"email":           config.StringVariable(email),
				"first_name":      config.StringVariable(firstName),
				"last_name":       config.StringVariable(lastName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created user
				resource.TestCheckResourceAttr("uptime_statuspage_user.test", "email", email),
				// Check that filtered output contains exactly 1 user
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
						return fmt.Errorf("expected exactly 1 filtered user, got %d", count)
					}
					return nil
				},
				// Check the filtered user has the correct email
				resource.TestCheckOutput("filtered_user_email", email),
			),
		},
	}))
}
