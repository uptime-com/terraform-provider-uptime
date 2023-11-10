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

func NewCheckICMPResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckICMPResourceModel, upapi.CheckICMP, upapi.Check]{
		CheckICMPResourceAPI{provider: p},
		CheckICMPResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_icmp",
			Schema: schema.Schema{
				Description: "Monitor network activity for a specific domain or IP address",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"address":                   AddressHostnameSchemaAttribute(),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"locations":                 LocationsSchemaAttribute(p),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"interval":                  IntervalSchemaAttribute(5),
					"num_retries":               NumRetriesAttribute(2),
					"use_ip_version":            UseIPVersionSchemaAttribute(),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
				},
			},
		},
	}
}

type CheckICMPResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	Address                types.String `tfsdk:"address"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Locations              types.Set    `tfsdk:"locations"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	UseIPVersion           types.String `tfsdk:"use_ip_version"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLA                    types.Object `tfsdk:"sla"`

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckICMPResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckICMPResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckICMPResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckICMPResourceModel, diag.Diagnostics) {
	model := *new(CheckICMPResourceModel)
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

func (a CheckICMPResourceModelAdapter) ToAPIArgument(model CheckICMPResourceModel) (*upapi.CheckICMP, error) {
	api := upapi.CheckICMP{
		Name:                   model.Name.ValueString(),
		Address:                model.Address.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		NumRetries:             model.NumRetries.ValueInt64(),
		UseIPVersion:           model.UseIPVersion.ValueString(),
		Notes:                  model.Notes.ValueString(),
		IncludeInGlobalMetrics: model.IncludeInGlobalMetrics.ValueBool(),
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

func (a CheckICMPResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckICMPResourceModel, error) {
	model := CheckICMPResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		Address:                types.StringValue(api.Address),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Locations:              a.LocationsValue(api.Locations),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		NumRetries:             types.Int64Value(api.NumRetries),
		UseIPVersion:           types.StringValue(api.UseIPVersion),
		Notes:                  types.StringValue(api.Notes),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Latency: DurationValueFromDecimalSeconds(api.ResponseTimeSLA),
			Uptime:  DecimalValue(api.UptimeSLA),
		}),
	}
	return &model, nil
}

type CheckICMPResourceAPI struct {
	provider *providerImpl
}

func (c CheckICMPResourceAPI) Create(ctx context.Context, arg upapi.CheckICMP) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateICMP(ctx, arg)
}

func (c CheckICMPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckICMPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckICMP) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateICMP(ctx, pk, arg)
}

func (c CheckICMPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
