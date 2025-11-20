package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPushNotificationProfilesDataSource(t *testing.T) {
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			Config: `data "uptime_push_notification_profiles" "test" {}`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.uptime_push_notification_profiles.test", "id", ""),
			),
		},
	}))
}
