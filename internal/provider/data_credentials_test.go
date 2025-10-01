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

func TestAccCredentialsDataSource(t *testing.T) {
	name := petname.Generate(3, "-")
	password := petname.Generate(1, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_credentials"),
			ConfigVariables: config.Variables{
				"display_name": config.StringVariable(name),
				"password":     config.StringVariable(password),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created credential
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", name),
				resource.TestCheckResourceAttr("uptime_credential.test", "credential_type", "BASIC"),
				// Check that filtered output contains exactly 1 credential
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
						return fmt.Errorf("expected exactly 1 filtered credential, got %d", count)
					}
					return nil
				},
				// Check the filtered credential has the correct name and type
				resource.TestCheckOutput("filtered_credential_name", name),
				resource.TestCheckOutput("filtered_credential_type", "BASIC"),
			),
		},
	}))
}
