package provider

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

var _ plancheck.PlanCheck = (*planCheckNoOp)(nil)

// planCheckNoOp is a plan check that does nothing. Just a debugging helper.
// E.g., adding this to a test case will allow you to set breakpoints and inspect plans and states.
//
//	     ...
//			resource.TestStep{
//			    ...
//				ConfigPlanChecks: resource.ConfigPlanChecks{
//					PreApply: []plancheck.PlanCheck{
//						&planCheckNoOp{},
//					},
//					PostApplyPreRefresh: []plancheck.PlanCheck{
//						&planCheckNoOp{},
//					},
//					PostApplyPostRefresh: []plancheck.PlanCheck{
//						&planCheckNoOp{},
//					},
//				},
//		       ...
//			...
type planCheckNoOp struct{}

func (c *planCheckNoOp) CheckPlan(ctx context.Context, rq plancheck.CheckPlanRequest, rs *plancheck.CheckPlanResponse) {
	return
}

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"uptime": providerserver.NewProtocol6WithError(VersionFactory("test")()),
	}
}

func testAccAPIClient(t testing.TB) upapi.API {
	t.Helper()

	token := os.Getenv("UPTIME_TOKEN")
	require.NotEmpty(t, token, "UPTIME_TOKEN must be set for acceptance tests")

	rateLimit := 0.15
	if val := os.Getenv("UPTIME_RATE_LIMIT"); val != "" {
		if parsedVal, err := strconv.ParseFloat(val, 64); err == nil {
			rateLimit = parsedVal
		}
	}

	opts := []upapi.Option{
		upapi.WithToken(token),
		upapi.WithRateLimit(rateLimit),
	}
	if endpoint := os.Getenv("UPTIME_ENDPOINT"); endpoint != "" {
		opts = append(opts, upapi.WithBaseURL(endpoint))
	}

	api, err := upapi.New(opts...)
	require.NoError(t, err, "failed to initialize uptime.com api client")

	return api
}

// testAccLocations fetches available probe server locations from the API.
// It returns a slice of location strings (excluding "AUTO" and "TEST").
// Skips the test if TF_ACC is not set.
func testAccLocations(t testing.TB) []string {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	api := testAccAPIClient(t)
	servers, err := api.ProbeServers().List(context.Background())
	require.NoError(t, err, "failed to fetch probe servers")

	locations := make([]string, 0, len(servers.Items))
	seen := make(map[string]struct{})
	for _, s := range servers.Items {
		if s.Location == "AUTO" || s.Location == "TEST" {
			continue
		}
		if _, ok := seen[s.Location]; ok {
			continue
		}
		seen[s.Location] = struct{}{}
		locations = append(locations, s.Location)
	}
	require.GreaterOrEqual(t, len(locations), 4, "expected at least 4 locations")
	return locations
}

func testCaseFromSteps(t testing.TB, steps []resource.TestStep) resource.TestCase {
	t.Helper()

	return resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps:                    steps,
	}
}
