package provider

import (
	"context"
	"testing"

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
	var matching []upapi.Check
	for _, c := range s.checks {
		if c.IsPaused == opts.IsPaused {
			matching = append(matching, c)
		}
	}
	start := (opts.Page - 1) * opts.PageSize
	end := min(start+opts.PageSize, int64(len(matching)))
	return &upapi.ListResult[upapi.Check]{Items: matching[start:end], TotalCount: int64(len(matching))}, nil
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

func checksForTest(active, paused int) []upapi.Check {
	var checks []upapi.Check
	for pk := int64(1); pk <= int64(active+paused); pk++ {
		checks = append(checks, upapi.Check{PK: pk, IsPaused: pk > int64(active)})
	}
	return checks
}

func TestCheckCache_ServesReadsFromList(t *testing.T) {
	checks := checksForTest(2*checkListPageSize+1, 1)
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
		{Page: 1, PageSize: checkListPageSize, IsPaused: true},
	}, ep.listCalls)
}

func TestCheckCache_ExactPageMultipleStopsWithoutExtraPage(t *testing.T) {
	cache, ep := newCheckCacheForTest(checksForTest(2*checkListPageSize, 0))

	_, err := cache.Get(context.Background(), upapi.PrimaryKey(1))

	require.NoError(t, err)
	assert.Len(t, ep.listCalls, 3, "2 active pages and 1 empty paused page")
}

func TestCheckCache_UnknownCheckFallsBackToGet(t *testing.T) {
	cache, ep := newCheckCacheForTest(checksForTest(1, 0))

	_, err := cache.Get(context.Background(), upapi.PrimaryKey(99))

	assert.Error(t, err)
	assert.Equal(t, []upapi.PrimaryKey{99}, ep.getCalls)
}

func TestCheckCache_ListFailureKeepsLoadedPagesAndFallsBackToGet(t *testing.T) {
	cache, ep := newCheckCacheForTest(checksForTest(checkListPageSize+1, 0))
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
	ep := &cacheStubChecks{checks: checksForTest(1, 0)}
	p := &providerImpl{api: cacheStubAPI{endpoint: ep}}

	_, err := p.getCheck(context.Background(), upapi.PrimaryKey(1))

	require.NoError(t, err)
	assert.Empty(t, ep.listCalls)
	assert.Equal(t, []upapi.PrimaryKey{1}, ep.getCalls)
}
