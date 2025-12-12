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

func NewCheckIMAPResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[CheckIMAPResourceModel, upapi.CheckIMAP, upapi.Check](
		CheckIMAPResourceAPI{provider: p},
		CheckIMAPResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_imap",
			Schema: schema.Schema{
				Description: "Monitor IMAP server availability. Import using the check ID: `terraform import uptime_check_imap.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"address":        AddressHostnameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"locations":      LocationsSchemaAttribute(p),
					"tags":           TagsSchemaAttribute(),
					"is_paused":      IsPausedSchemaAttribute(),
					"interval":       IntervalSchemaAttribute(5),
					"port":           PortSchemaAttribute(143),
					"expect_string":  StringToExpectSchemaAttribute(),
					"encryption": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Whether to use TLS",
					},
					"use_ip_version":            UseIPVersionSchemaAttribute(),
					"num_retries":               NumRetriesAttribute(2),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
				},
			},
		},
		ImportStateSimpleID,
	)
}

type CheckIMAPResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	Address                types.String `tfsdk:"address"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Locations              types.Set    `tfsdk:"locations"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	Port                   types.Int64  `tfsdk:"port"`
	ExpectString           types.String `tfsdk:"expect_string"`
	Encryption             types.String `tfsdk:"encryption"`
	UseIPVersion           types.String `tfsdk:"use_ip_version"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLA                    types.Object `tfsdk:"sla"`

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckIMAPResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckIMAPResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckIMAPResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckIMAPResourceModel, diag.Diagnostics) {
	model := *new(CheckIMAPResourceModel)
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

func (a CheckIMAPResourceModelAdapter) ToAPIArgument(model CheckIMAPResourceModel) (*upapi.CheckIMAP, error) {
	api := upapi.CheckIMAP{
		Name:                   model.Name.ValueString(),
		Address:                model.Address.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		Port:                   model.Port.ValueInt64(),
		ExpectString:           model.ExpectString.ValueString(),
		Encryption:             model.Encryption.ValueString(),
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

func (a CheckIMAPResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckIMAPResourceModel, error) {
	model := CheckIMAPResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		Address:                types.StringValue(api.Address),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Locations:              a.LocationsValue(api.Locations),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		Port:                   types.Int64Value(api.Port),
		ExpectString:           types.StringValue(api.ExpectString),
		Encryption:             types.StringValue(api.Encryption),
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

type CheckIMAPResourceAPI struct {
	provider *providerImpl
}

func (c CheckIMAPResourceAPI) Create(ctx context.Context, arg upapi.CheckIMAP) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateIMAP(ctx, arg)
}

func (c CheckIMAPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckIMAPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckIMAP) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateIMAP(ctx, pk, arg)
}

func (c CheckIMAPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
