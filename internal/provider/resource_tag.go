package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewTagResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[TagResourceModel, upapi.Tag, upapi.Tag]{
		TagResourceAPI{provider: p},
		TagResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "tag",
			Schema:         TagResourceSchema,
		},
	}
}

var TagResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id":  IDSchemaAttribute(),
		"url": URLSchemaAttribute(),
		"tag": schema.StringAttribute{
			Required: true,
		},
		"color_hex": ColorHexSchemaAttribute(),
	},
}

type TagResourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	URL      types.String `tfsdk:"url"`
	Tag      types.String `tfsdk:"tag"`
	ColorHex types.String `tfsdk:"color_hex"`
}

func (m TagResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type TagResourceModelAdapter struct{}

func (TagResourceModelAdapter) Get(ctx context.Context, getter StateGetter) (*TagResourceModel, diag.Diagnostics) {
	var m TagResourceModel
	diags := getter.Get(ctx, &m)
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (TagResourceModelAdapter) ToAPIArgument(m TagResourceModel) (*upapi.Tag, error) {
	return &upapi.Tag{
		Tag:      m.Tag.ValueString(),
		ColorHex: m.ColorHex.ValueString(),
	}, nil
}

func (TagResourceModelAdapter) FromAPIResult(obj upapi.Tag) (*TagResourceModel, error) {
	return &TagResourceModel{
		ID:       types.Int64Value(obj.PK),
		URL:      types.StringValue(obj.URL),
		Tag:      types.StringValue(obj.Tag),
		ColorHex: types.StringValue(obj.ColorHex),
	}, nil
}

type TagResourceAPI struct {
	provider *providerImpl
}

func (c TagResourceAPI) Create(ctx context.Context, arg upapi.Tag) (*upapi.Tag, error) {
	obj, err := c.provider.api.Tags().Create(ctx, arg)
	return obj, err
}

func (c TagResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Tag, error) {
	obj, err := c.provider.api.Tags().Get(ctx, pk)
	return obj, err
}

func (c TagResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.Tag) (*upapi.Tag, error) {
	obj, err := c.provider.api.Tags().Update(ctx, pk, arg)
	return obj, err
}

func (c TagResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Tags().Delete(ctx, pk)
}
