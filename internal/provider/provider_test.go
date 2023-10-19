package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"uptime": providerserver.NewProtocol6WithError(VersionFactory("test")()),
	}
}

func testAccAPIClient(t testing.TB) upapi.API {
	t.Helper()

	token := os.Getenv("UPTIME_TOKEN")
	require.NotEmpty(t, token, "UPTIME_TOKEN must be set for acceptance tests")

	api, err := upapi.New(upapi.WithToken(token), upapi.WithRateLimit(0.2))
	require.NoError(t, err, "failed to initialize uptime.com api client")

	return api
}

func testCaseFromSteps(t testing.TB, steps []resource.TestStep) resource.TestCase {
	t.Helper()

	return resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps:                    steps,
	}
}
