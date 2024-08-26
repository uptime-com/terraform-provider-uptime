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

func NewCheckTCPResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckTCPResourceModel, upapi.CheckTCP, upapi.Check]{
		CheckTCPResourceAPI{provider: p},
		CheckTCPResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_tcp",
			Schema: schema.Schema{
				Description: "Monitor a TCP port for a response",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"address":                   AddressHostnameSchemaAttribute(),
					"port":                      RequiredPortSchemaAttribute(),
					"send_string":               StringToSendSchemaAttribute(),
					"expect_string":             StringToExpectSchemaAttribute(),
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
					"encryption": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Whether to use TLS",
					},
				},
			},
		},
	}
}

type CheckTCPResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	Address                types.String `tfsdk:"address"`
	Port                   types.Int64  `tfsdk:"port"`
	SendString             types.String `tfsdk:"send_string"`
	ExpectString           types.String `tfsdk:"expect_string"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Locations              types.Set    `tfsdk:"locations"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	UseIpVersion           types.String `tfsdk:"use_ip_version"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLA                    types.Object `tfsdk:"sla"`
	Encryption             types.String `tfsdk:"encryption"`

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckTCPResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckTCPResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckTCPResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckTCPResourceModel, diag.Diagnostics) {
	model := *new(CheckTCPResourceModel)
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

func (a CheckTCPResourceModelAdapter) ToAPIArgument(model CheckTCPResourceModel) (*upapi.CheckTCP, error) {
	api := upapi.CheckTCP{
		Name:                   model.Name.ValueString(),
		Address:                model.Address.ValueString(),
		Port:                   model.Port.ValueInt64(),
		SendString:             model.SendString.ValueString(),
		ExpectString:           model.ExpectString.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		NumRetries:             model.NumRetries.ValueInt64(),
		UseIpVersion:           model.UseIpVersion.ValueString(),
		Notes:                  model.Notes.ValueString(),
		IncludeInGlobalMetrics: model.IncludeInGlobalMetrics.ValueBool(),
		Encryption:             model.Encryption.ValueString(),
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

func (a CheckTCPResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckTCPResourceModel, error) {
	model := CheckTCPResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		Address:                types.StringValue(api.Address),
		Port:                   types.Int64Value(api.Port),
		SendString:             types.StringValue(api.SendString),
		ExpectString:           types.StringValue(api.ExpectString),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Locations:              a.LocationsValue(api.Locations),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		NumRetries:             types.Int64Value(api.NumRetries),
		UseIpVersion:           types.StringValue(api.UseIPVersion),
		Notes:                  types.StringValue(api.Notes),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Latency: DurationValueFromDecimalSeconds(api.ResponseTimeSLA),
			Uptime:  DecimalValue(api.UptimeSLA),
		}),
		Encryption: types.StringValue(api.Encryption),
	}
	return &model, nil
}

type CheckTCPResourceAPI struct {
	provider *providerImpl
}

func (m *CheckTCPResourceModel) APIResource() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

func (c CheckTCPResourceAPI) Create(ctx context.Context, arg upapi.CheckTCP) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateTCP(ctx, arg)
}

func (c CheckTCPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckTCPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckTCP) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateTCP(ctx, pk, arg)
}

func (c CheckTCPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
