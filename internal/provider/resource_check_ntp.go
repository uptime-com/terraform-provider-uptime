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

func NewCheckNTPResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckNTPResourceModel, upapi.CheckNTP, upapi.Check]{
		api: CheckNTPResourceAPI{provider: p},
		mod: CheckNTPResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_ntp",
			Schema: schema.Schema{
				Description: "Monitor a Network Time Protocol server",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"address":                   AddressHostnameSchemaAttribute(),
					"port":                      PortSchemaAttribute(123),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"locations":                 LocationsSchemaAttribute(p),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"interval":                  IntervalSchemaAttribute(5),
					"threshold":                 ThresholdSchemaAttribute(20),
					"sensitivity":               SensitivitySchemaAttribute(2),
					"num_retries":               NumRetriesSchemaAttribute(2),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"use_ip_version":            UseIPVersionSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
				},
			},
		},
	}
}

// var checkNTPResourceSchema =
type CheckNTPResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	URL                    types.String `tfsdk:"url"`
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
	SLA                    types.Object `tfsdk:"sla"`

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckNTPResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckNTPResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckNTPResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckNTPResourceModel, diag.Diagnostics) {
	model := *new(CheckNTPResourceModel)
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

func (a CheckNTPResourceModelAdapter) ToAPIArgument(model CheckNTPResourceModel) (*upapi.CheckNTP, error) {
	api := upapi.CheckNTP{
		Name:                   model.Name.ValueString(),
		Address:                model.Address.ValueString(),
		Port:                   model.Port.ValueInt64(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		Threshold:              model.Threshold.ValueInt64(),
		Sensitivity:            model.Sensitivity.ValueInt64(),
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

func (a CheckNTPResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckNTPResourceModel, error) {
	model := CheckNTPResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		Address:                types.StringValue(api.Address),
		Port:                   types.Int64Value(api.Port),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Locations:              a.LocationsValue(api.Locations),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		Threshold:              types.Int64Value(api.Threshold),
		Sensitivity:            types.Int64Value(api.Sensitivity),
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

type CheckNTPResourceAPI struct {
	provider *providerImpl
}

func (c CheckNTPResourceAPI) Create(ctx context.Context, arg upapi.CheckNTP) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateNTP(ctx, arg)
}

func (c CheckNTPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckNTPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckNTP) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateNTP(ctx, pk, arg)
}

func (c CheckNTPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
