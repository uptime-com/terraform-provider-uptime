package provider

import (
	"context"
	"sync/atomic"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

type cacheStubAPI struct {
	upapi.API
	endpoint *cacheStubChecks
}

func (s cacheStubAPI) Checks() upapi.ChecksEndpoint { return s.endpoint }

type cacheStubChecks struct {
	upapi.ChecksEndpoint
	checks    []upapi.Check
	listCalls []upapi.CheckListOptions
	getCalls  []upapi.PrimaryKey
	failPage  int64
}

func (s *cacheStubChecks) List(_ context.Context, opts upapi.CheckListOptions) (*upapi.ListResult[upapi.Check], error) {
	s.listCalls = append(s.listCalls, opts)
	if opts.Page == s.failPage {
		return nil, assert.AnError
	}
	start := (opts.Page - 1) * opts.PageSize
	end := min(start+opts.PageSize, int64(len(s.checks)))
	return &upapi.ListResult[upapi.Check]{Items: s.checks[start:end], TotalCount: int64(len(s.checks))}, nil
}

func (s *cacheStubChecks) Get(_ context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	s.getCalls = append(s.getCalls, pk.PrimaryKey())
	for _, c := range s.checks {
		if c.PrimaryKey() == pk.PrimaryKey() {
			return &c, nil
		}
	}
	return nil, &upapi.Error{}
}

func newCheckCacheForTest(checks []upapi.Check) (*checkCache, *cacheStubChecks) {
	ep := &cacheStubChecks{checks: checks}
	return &checkCache{api: cacheStubAPI{endpoint: ep}}, ep
}

func checksForTest(n int) []upapi.Check {
	var checks []upapi.Check
	for pk := int64(1); pk <= int64(n); pk++ {
		checks = append(checks, upapi.Check{PK: pk, IsPaused: pk%2 == 0})
	}
	return checks
}

func TestCheckCache_ServesReadsFromList(t *testing.T) {
	checks := checksForTest(2*checkListPageSize + 1)
	cache, ep := newCheckCacheForTest(checks)

	for _, c := range checks {
		got, err := cache.Get(context.Background(), upapi.PrimaryKey(c.PK))
		require.NoError(t, err)
		assert.Equal(t, c.PK, got.PK)
	}

	assert.Empty(t, ep.getCalls, "every check must come from the list")
	assert.Equal(t, []upapi.CheckListOptions{
		{Page: 1, PageSize: checkListPageSize},
		{Page: 2, PageSize: checkListPageSize},
		{Page: 3, PageSize: checkListPageSize},
	}, ep.listCalls)
}

func TestCheckCache_ExactPageMultipleStopsWithoutExtraPage(t *testing.T) {
	cache, ep := newCheckCacheForTest(checksForTest(2 * checkListPageSize))

	_, err := cache.Get(context.Background(), upapi.PrimaryKey(1))

	require.NoError(t, err)
	assert.Len(t, ep.listCalls, 2)
}

func TestCheckCache_UnknownCheckFallsBackToGet(t *testing.T) {
	cache, ep := newCheckCacheForTest(checksForTest(1))

	_, err := cache.Get(context.Background(), upapi.PrimaryKey(99))

	assert.Error(t, err)
	assert.Equal(t, []upapi.PrimaryKey{99}, ep.getCalls)
}

func TestCheckCache_ListFailureKeepsLoadedPagesAndFallsBackToGet(t *testing.T) {
	cache, ep := newCheckCacheForTest(checksForTest(checkListPageSize + 1))
	ep.failPage = 2

	first, err := cache.Get(context.Background(), upapi.PrimaryKey(1))
	require.NoError(t, err)
	last, err := cache.Get(context.Background(), upapi.PrimaryKey(checkListPageSize+1))
	require.NoError(t, err)

	assert.EqualValues(t, 1, first.PK)
	assert.EqualValues(t, checkListPageSize+1, last.PK)
	assert.Len(t, ep.listCalls, 2, "a failed load is not retried")
	assert.Equal(t, []upapi.PrimaryKey{checkListPageSize + 1}, ep.getCalls)
}

func TestProviderGetCheck_UsesGetWithoutBulkRead(t *testing.T) {
	ep := &cacheStubChecks{checks: checksForTest(1)}
	p := &providerImpl{api: cacheStubAPI{endpoint: ep}}

	_, err := p.getCheck(context.Background(), upapi.PrimaryKey(1))

	require.NoError(t, err)
	assert.Empty(t, ep.listCalls)
	assert.Equal(t, []upapi.PrimaryKey{1}, ep.getCalls)
}

type countingAPI struct {
	upapi.API
	lists, gets *atomic.Int64
}

func (a countingAPI) Checks() upapi.ChecksEndpoint {
	return countingChecks{ChecksEndpoint: a.API.Checks(), lists: a.lists, gets: a.gets}
}

type countingChecks struct {
	upapi.ChecksEndpoint
	lists, gets *atomic.Int64
}

func (c countingChecks) List(ctx context.Context, opts upapi.CheckListOptions) (*upapi.ListResult[upapi.Check], error) {
	c.lists.Add(1)
	return c.ChecksEndpoint.List(ctx, opts)
}

func (c countingChecks) Get(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	c.gets.Add(1)
	return c.ChecksEndpoint.Get(ctx, pk)
}

// The test framework starts a new provider server for every Terraform command, so a fresh
// providerImpl per factory call gives the cache the same one-walk lifetime as a real run.
func TestAccBulkRead_CheckHTTP(t *testing.T) {
	var lists, gets atomic.Int64
	factories := map[string]func() (tfprotov6.ProviderServer, error){
		"uptime": func() (tfprotov6.ProviderServer, error) {
			api := countingAPI{API: testAccAPIClient(t), lists: &lists, gets: &gets}
			p := &providerImpl{version: "test", api: api, checks: &checkCache{api: api}}
			return providerserver.NewProtocol6WithError(p)()
		},
	}
	steps := make([]resource.TestStep, 0, 2)
	for _, header := range []string{"Bar", "Baz"} {
		name := petname.Generate(3, "-")
		steps = append(steps, resource.TestStep{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/headers"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"headers": config.MapVariable(map[string]config.Variable{
					"Foo": config.ListVariable(config.StringVariable(header)),
				}),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "name", name),
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Foo.0", header),
			),
		})
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: factories,
		Steps:                    steps,
	})

	assert.Positive(t, lists.Load(), "checks must be loaded from the list endpoint")
	assert.Zero(t, gets.Load(), "no per-check GET must be needed for checks the list returned")
}
