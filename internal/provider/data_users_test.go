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

func TestAccUsersDataSource(t *testing.T) {
	firstName := petname.Generate(2, "-")
	lastName := petname.Generate(2, "-")
	email := fmt.Sprintf("%s@example.com", petname.Generate(3, "-"))
	password := fmt.Sprintf("%s123", petname.Generate(4, "-"))

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_users"),
			ConfigVariables: config.Variables{
				"first_name": config.StringVariable(firstName),
				"last_name":  config.StringVariable(lastName),
				"email":      config.StringVariable(email),
				"password":   config.StringVariable(password),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created user
				resource.TestCheckResourceAttr("uptime_user.test", "first_name", firstName),
				resource.TestCheckResourceAttr("uptime_user.test", "email", email),
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
