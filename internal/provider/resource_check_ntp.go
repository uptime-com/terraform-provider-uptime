package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckNTPResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkNTPResourceModel, upapi.CheckNTP, upapi.Check]{
		api: &checkNTPResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_ntp",
			Schema:         checkNTPResourceSchema,
		},
	}
}

var checkNTPResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Optional: true,
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
		"threshold": schema.Int64Attribute{
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
		"use_ip_version": schema.StringAttribute{
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

type checkNTPResourceModel struct {
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
	Threshold              types.Int64  `tfsdk:"threshold"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	UseIPVersion           types.String `tfsdk:"use_ip_version"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
}

var _ genericResourceAPI[upapi.CheckNTP, upapi.Check] = (*checkNTPResourceAPI)(nil)

type checkNTPResourceAPI struct {
	provider *providerImpl
}

func (c *checkNTPResourceAPI) Create(ctx context.Context, arg upapi.CheckNTP) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateNTP(ctx, arg)
}

func (c *checkNTPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c *checkNTPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckNTP) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateNTP(ctx, pk, arg)
}

func (c *checkNTPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
