package provider

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPushNotificationsDataSource(t *testing.T) {
	// Push notification tests require actual mobile device credentials.
	// Unlike other integrations (e.g., Opsgenie), the Uptime.com API validates
	// push notification credentials at creation time, rejecting invalid app_keys.
	// Set UPTIME_TEST_PUSH_APP_KEY and UPTIME_TEST_PUSH_UUID from a real device.
	appKey := os.Getenv("UPTIME_TEST_PUSH_APP_KEY")
	uuid := os.Getenv("UPTIME_TEST_PUSH_UUID")

	if appKey == "" || uuid == "" {
		t.Skip("Skipping push notifications data source test: UPTIME_TEST_PUSH_APP_KEY and UPTIME_TEST_PUSH_UUID must be set with valid mobile device credentials from the Uptime.com mobile app")
	}

	deviceName := petname.Generate(3, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/data_push_notifications"),
			ConfigVariables: config.Variables{
				"device_name": config.StringVariable(deviceName),
				"uuid":        config.StringVariable(uuid),
				"app_key":     config.StringVariable(appKey),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check the created push notification profile
				resource.TestCheckResourceAttr("uptime_push_notification.test", "device_name", deviceName),
				resource.TestCheckResourceAttr("uptime_push_notification.test", "uuid", uuid),
				// Check that filtered output contains exactly 1 profile
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
						return fmt.Errorf("expected exactly 1 filtered push notification profile, got %d", count)
					}
					return nil
				},
				// Check the filtered profile has the correct device name
				resource.TestCheckOutput("filtered_device_name", deviceName),
				resource.TestCheckOutput("filtered_uuid", uuid),
			),
		},
	}))
}
