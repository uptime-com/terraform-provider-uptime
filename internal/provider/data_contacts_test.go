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

func TestAccContactsDataSource(t *testing.T) {
	name := petname.Generate(3, "-")
	email := fmt.Sprintf("%s@example.com", petname.Generate(2, "-"))

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_contacts"),
			ConfigVariables: config.Variables{
				"name":  config.StringVariable(name),
				"email": config.StringVariable(email),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created contact
				resource.TestCheckResourceAttr("uptime_contact.test", "name", name),
				// Check that filtered output contains exactly 1 contact
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
						return fmt.Errorf("expected exactly 1 filtered contact, got %d", count)
					}
					return nil
				},
				// Check the filtered contact has the correct name
				resource.TestCheckOutput("filtered_contact_name", name),
			),
		},
	}))
}
