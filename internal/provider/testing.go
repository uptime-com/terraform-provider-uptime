package provider

import (
	"context"
	"os"
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

	api, err := upapi.New(upapi.WithToken(token), upapi.WithRateLimit(0.15))
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
