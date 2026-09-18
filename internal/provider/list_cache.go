package provider

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

const listPageSize = 250

// listCache serves resource Reads from a single paginated List instead of one Get per item.
// The first call loads every page once; an item missing from the cache falls back to Get so
// 404 handling stays unchanged. The snapshot lives for one provider Configure, which Terraform
// runs once per plan or apply walk, and those walks never mix Read with Create, Update or Delete.
type listCache[T upapi.PrimaryKeyable] struct {
	name  string
	list  func(ctx context.Context, page int64) (*upapi.ListResult[T], error)
	get   func(ctx context.Context, pk upapi.PrimaryKeyable) (*T, error)
	once  sync.Once
	items map[upapi.PrimaryKey]T
}

func newCheckCache(api upapi.API) *listCache[upapi.Check] {
	checks := api.Checks()
	return &listCache[upapi.Check]{
		name: "checks",
		list: func(ctx context.Context, page int64) (*upapi.ListResult[upapi.Check], error) {
			return checks.List(ctx, upapi.CheckListOptions{Page: page, PageSize: listPageSize})
		},
		get: checks.Get,
	}
}

func newTagCache(api upapi.API) *listCache[upapi.Tag] {
	tags := api.Tags()
	return &listCache[upapi.Tag]{
		name: "tags",
		list: func(ctx context.Context, page int64) (*upapi.ListResult[upapi.Tag], error) {
			return tags.List(ctx, upapi.TagListOptions{Page: page, PageSize: listPageSize})
		},
		get: tags.Get,
	}
}

func (c *listCache[T]) load(ctx context.Context) map[upapi.PrimaryKey]T {
	items := make(map[upapi.PrimaryKey]T)
	for page := int64(1); ; page++ {
		res, err := c.list(ctx, page)
		if err != nil {
			tflog.Warn(ctx, "bulk_read: list failed, missing items are read one by one",
				map[string]any{"resource": c.name, "error": err.Error()})
			return items
		}
		for _, item := range res.Items {
			items[item.PrimaryKey()] = item
		}
		if page*listPageSize >= res.TotalCount {
			return items
		}
	}
}

func (c *listCache[T]) Get(ctx context.Context, pk upapi.PrimaryKeyable) (*T, error) {
	c.once.Do(func() { c.items = c.load(ctx) })
	if item, ok := c.items[pk.PrimaryKey()]; ok {
		return &item, nil
	}
	return c.get(ctx, pk)
}
