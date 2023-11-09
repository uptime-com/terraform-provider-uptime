package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shopspring/decimal"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckHeartbeatResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckHeartbeatResourceModel, upapi.CheckHeartbeat, upapi.Check]{
		api: CheckHeartbeatResourceAPI{provider: p},
		mod: CheckHeartbeatResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_heartbeat",
			Schema: schema.Schema{
				Description: "Monitor a periodic process, such as Cron, and issue alerts if the expected interval is exceeded",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"interval":                  IntervalSchemaAttribute(5),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"heartbeat_url": schema.StringAttribute{
						Computed:    true,
						Description: "URL to send data to the check",
					},
					"sla": SLASchemaAttribute(),
				},
			},
		},
	}
}

type CheckHeartbeatResourceModel struct {
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
	SLA                    types.Object `tfsdk:"sla"`

	sla *SLAAttribute `tfsdk:"-"`
}

type CheckHeartbeatResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckHeartbeatResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckHeartbeatResourceModel, diag.Diagnostics) {
	model := *new(CheckHeartbeatResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.sla, diags = a.SLAAttributeContext(ctx, model.SLA)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckHeartbeatResourceModelAdapter) ToAPIArgument(model CheckHeartbeatResourceModel) (_ *upapi.CheckHeartbeat, err error) {
	api := upapi.CheckHeartbeat{
		Name:                   model.Name.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		Notes:                  model.Notes.ValueString(),
		IncludeInGlobalMetrics: model.IncludeInGlobalMetrics.ValueBool(),
		HeartbeatURL:           model.HeartbeatURL.ValueString(),
	}

	if model.sla != nil {
		if !model.sla.Uptime.IsUnknown() {
			api.UptimeSLA = model.sla.Uptime.ValueDecimal()
		}
		if !model.sla.Latency.IsUnknown() {
			api.ResponseTimeSLA = decimal.NewFromFloat(model.sla.Latency.ValueDuration().Seconds())
		}
	}

	return &api, nil
}

func (a CheckHeartbeatResourceModelAdapter) FromAPIResult(api upapi.Check) (_ *CheckHeartbeatResourceModel, err error) {
	model := CheckHeartbeatResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		Notes:                  types.StringValue(api.Notes),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
		HeartbeatURL:           types.StringValue(api.HeartbeatURL),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Uptime:  DecimalValue(api.UptimeSLA),
			Latency: DurationValueFromDecimalSeconds(api.ResponseTimeSLA),
		}),
	}
	return &model, nil
}

func (m CheckHeartbeatResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckHeartbeatResourceAPI struct {
	provider *providerImpl
}

func (a CheckHeartbeatResourceAPI) Create(ctx context.Context, arg upapi.CheckHeartbeat) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateHeartbeat(ctx, arg)
}

func (a CheckHeartbeatResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a CheckHeartbeatResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckHeartbeat) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateHeartbeat(ctx, pk, arg)
}

func (a CheckHeartbeatResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
