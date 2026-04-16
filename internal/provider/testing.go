package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

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

var (
	testAccAPIClientOnce sync.Once
	testAccAPIClientInst upapi.API
	testAccAPIClientErr  error

	testAccProviderOnce sync.Once
	testAccProviderInst *providerImpl
	testAccProviderErr  error
)

func buildTestAccAPIClient() (upapi.API, error) {
	token := os.Getenv("UPTIME_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("UPTIME_TOKEN must be set for acceptance tests")
	}
	rateLimit := 0.5
	if val := os.Getenv("UPTIME_RATE_LIMIT"); val != "" {
		if parsedVal, err := strconv.ParseFloat(val, 64); err == nil {
			rateLimit = parsedVal
		}
	}
	var subaccount int64
	if val := os.Getenv("UPTIME_SUBACCOUNT"); val != "" {
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UPTIME_SUBACCOUNT: %w", err)
		}
		subaccount = parsed
	}
	opts := []upapi.Option{
		upapi.WithSubaccount(subaccount),
		upapi.WithToken(token),
		upapi.WithUserAgent((&providerImpl{version: "test"}).UserAgentString()),
		upapi.WithRateLimit(rateLimit),
		upapi.WithRetry(10, time.Second*30, os.Stderr),
	}
	if endpoint := os.Getenv("UPTIME_ENDPOINT"); endpoint != "" {
		opts = append(opts, upapi.WithBaseURL(endpoint))
	}
	if os.Getenv("UPTIME_TRACE") != "" {
		opts = append(opts, upapi.WithTrace(os.Stderr))
	}
	return upapi.New(opts...)
}

func sharedTestAccAPIClient() (upapi.API, error) {
	testAccAPIClientOnce.Do(func() {
		testAccAPIClientInst, testAccAPIClientErr = buildTestAccAPIClient()
	})
	return testAccAPIClientInst, testAccAPIClientErr
}

func sharedTestAccProvider() (*providerImpl, error) {
	testAccProviderOnce.Do(func() {
		api, err := sharedTestAccAPIClient()
		if err != nil {
			testAccProviderErr = err
			return
		}
		p, ok := VersionFactory("test")().(*providerImpl)
		if !ok {
			testAccProviderErr = fmt.Errorf("VersionFactory did not return *providerImpl")
			return
		}
		p.api = api
		testAccProviderInst = p
	})
	return testAccProviderInst, testAccProviderErr
}

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"uptime": func() (tfprotov6.ProviderServer, error) {
			p, err := sharedTestAccProvider()
			if err != nil {
				return nil, err
			}
			return providerserver.NewProtocol6WithError(p)()
		},
	}
}

func testAccAPIClient(t testing.TB) upapi.API {
	t.Helper()

	api, err := sharedTestAccAPIClient()
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
