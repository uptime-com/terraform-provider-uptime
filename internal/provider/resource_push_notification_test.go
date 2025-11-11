package provider

import (
	"os"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPushNotificationResource(t *testing.T) {
	// Push notification tests require actual mobile device credentials.
	// Unlike other integrations (e.g., Opsgenie), the Uptime.com API validates
	// push notification credentials at creation time, rejecting invalid app_keys.
	// Set UPTIME_TEST_PUSH_APP_KEY and UPTIME_TEST_PUSH_UUID from a real device.
	appKey := os.Getenv("UPTIME_TEST_PUSH_APP_KEY")
	uuid := os.Getenv("UPTIME_TEST_PUSH_UUID")

	if appKey == "" || uuid == "" {
		t.Skip("Skipping push notification test: UPTIME_TEST_PUSH_APP_KEY and UPTIME_TEST_PUSH_UUID must be set with valid mobile device credentials from the Uptime.com mobile app")
	}

	deviceNames := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}

	// Use the same credentials for both test steps (changing device name only)
	uuids := [2]string{uuid, uuid}
	appKeys := [2]string{appKey, appKey}

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"device_name": config.StringVariable(deviceNames[0]),
				"uuid":        config.StringVariable(uuids[0]),
				"app_key":     config.StringVariable(appKeys[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_push_notification/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_push_notification.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_push_notification.test", "url"),
				resource.TestCheckResourceAttr("uptime_push_notification.test", "device_name", deviceNames[0]),
				resource.TestCheckResourceAttr("uptime_push_notification.test", "uuid", uuids[0]),
				resource.TestCheckResourceAttr("uptime_push_notification.test", "app_key", appKeys[0]),
				resource.TestCheckResourceAttrSet("uptime_push_notification.test", "display_name"),
				resource.TestCheckResourceAttrSet("uptime_push_notification.test", "user"),
				resource.TestCheckResourceAttrSet("uptime_push_notification.test", "created_at"),
				resource.TestCheckResourceAttrSet("uptime_push_notification.test", "modified_at"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"device_name": config.StringVariable(deviceNames[1]),
				"uuid":        config.StringVariable(uuids[1]),
				"app_key":     config.StringVariable(appKeys[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_push_notification/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_push_notification.test", "device_name", deviceNames[1]),
				resource.TestCheckResourceAttr("uptime_push_notification.test", "uuid", uuids[1]),
				resource.TestCheckResourceAttr("uptime_push_notification.test", "app_key", appKeys[1]),
			),
		},
	}))
}
