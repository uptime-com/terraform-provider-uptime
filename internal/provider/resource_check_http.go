package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckHTTPResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkHTTPResourceModel, upapi.CheckHTTP, upapi.Check]{
		api: &checkHTTPResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_http",
			Schema:         checkHTTPResourceSchema,
		},
	}
}

var checkHTTPResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"contact_groups": schema.SetAttribute{
			ElementType: types.StringType,
			Required:    true,
		},
		"locations": schema.SetAttribute{
			ElementType: types.StringType,
			Required:    true,
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"is_paused": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"interval": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"address": schema.StringAttribute{
			Required: true,
		},
		"port": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"username": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"password": schema.StringAttribute{
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"proxy": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"status_code": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"send_string": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"expect_string": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"expect_string_type": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"encryption": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"threshold": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"headers": schema.MapAttribute{
			ElementType: types.ListType{
				ElemType: types.StringType,
			},
			Optional: true,
			Computed: true,
		},
		"version": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"sensitivity": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"num_retries": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"notes": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"include_in_global_metrics": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
	},
}

type checkHTTPResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Locations              types.Set    `tfsdk:"locations"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	Address                types.String `tfsdk:"address"`
	Port                   types.Int64  `tfsdk:"port"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	Proxy                  types.String `tfsdk:"proxy"`
	StatusCode             types.String `tfsdk:"status_code"`
	SendString             types.String `tfsdk:"send_string"`
	ExpectString           types.String `tfsdk:"expect_string"`
	ExpectStringType       types.String `tfsdk:"expect_string_type"`
	Encryption             types.String `tfsdk:"encryption"`
	Threshold              types.Int64  `tfsdk:"threshold"`
	Headers                types.Map    `tfsdk:"headers" ref:",extra=headers"`
	Version                types.Int64  `tfsdk:"version"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
}

var _ genericResourceAPI[upapi.CheckHTTP, upapi.Check] = (*checkHTTPResourceAPI)(nil)

type checkHTTPResourceAPI struct {
	provider *providerImpl
}

func (a *checkHTTPResourceAPI) Create(ctx context.Context, arg upapi.CheckHTTP) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateHTTP(ctx, arg)
}

func (a *checkHTTPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a *checkHTTPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckHTTP) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateHTTP(ctx, pk, arg)
}

func (a *checkHTTPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
