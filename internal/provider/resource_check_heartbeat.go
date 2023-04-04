package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckHeartbeatResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkHeartbeatResourceModel, upapi.CheckHeartbeat, upapi.Check]{
		api: &checkHeartbeatResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_heartbeat",
			Schema:         checkHeartbeatResourceSchema,
		},
	}
}

var checkHeartbeatResourceSchema = schema.Schema{
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
			Required:    true,
			ElementType: types.StringType,
		},
		"tags": schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
		},
		"is_paused": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"interval": schema.Int64Attribute{
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
		"heartbeat_url": schema.StringAttribute{
			Computed: true,
		},
	},
}

type checkHeartbeatResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"   ref:"PK,opt"`
	URL                    types.String `tfsdk:"url"  ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	HeartbeatURL           types.String `tfsdk:"heartbeat_url"`
}

var _ genericResourceAPI[upapi.CheckHeartbeat, upapi.Check] = (*checkHeartbeatResourceAPI)(nil)

type checkHeartbeatResourceAPI struct {
	provider *providerImpl
}

func (a *checkHeartbeatResourceAPI) Create(ctx context.Context, arg upapi.CheckHeartbeat) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateHeartbeat(ctx, arg)
}

func (a *checkHeartbeatResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a *checkHeartbeatResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckHeartbeat) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateHeartbeat(ctx, pk, arg)
}

func (a *checkHeartbeatResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
