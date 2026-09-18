package provider

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

const checkListPageSize = 250

// checkCache serves check Reads from a single paginated List instead of one Get per check.
// The first call loads every page once; a check missing from the cache falls back to Get so
// 404 handling stays unchanged. The snapshot lives for one provider Configure, which Terraform
// runs once per plan or apply walk, and those walks never mix Read with Create, Update or Delete.
type checkCache struct {
	api   upapi.API
	once  sync.Once
	items map[upapi.PrimaryKey]upapi.Check
}

func (c *checkCache) load(ctx context.Context) map[upapi.PrimaryKey]upapi.Check {
	items := make(map[upapi.PrimaryKey]upapi.Check)
	for page := int64(1); ; page++ {
		res, err := c.api.Checks().List(ctx, upapi.CheckListOptions{Page: page, PageSize: checkListPageSize})
		if err != nil {
			tflog.Warn(ctx, "bulk_read: listing checks failed, missing checks are read one by one",
				map[string]any{"error": err.Error()})
			return items
		}
		for _, item := range res.Items {
			items[item.PrimaryKey()] = item
		}
		if page*checkListPageSize >= res.TotalCount {
			return items
		}
	}
}

func (c *checkCache) Get(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	c.once.Do(func() { c.items = c.load(ctx) })
	if item, ok := c.items[pk.PrimaryKey()]; ok {
		return &item, nil
	}
	return c.api.Checks().Get(ctx, pk)
}
