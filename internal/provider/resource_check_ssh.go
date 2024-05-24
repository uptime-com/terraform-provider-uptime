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

func NewCheckSSHResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckSSHResourceModel, upapi.CheckSSH, upapi.Check]{
		CheckSSHResourceAPI{provider: p},
		CheckSSHResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_ssh",
			Schema: schema.Schema{
				Description: "Monitor SSH access for a domain or IP address",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"address":                   AddressHostnameSchemaAttribute(),
					"port":                      RequiredPortSchemaAttribute(),
					"sensitivity":               SensitivitySchemaAttribute(2),
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

type CheckSSHResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	Address                types.String `tfsdk:"address"`
	Port                   types.Int64  `tfsdk:"port"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"`
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

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckSSHResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckSSHResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckSSHResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckSSHResourceModel, diag.Diagnostics) {
	model := *new(CheckSSHResourceModel)
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

func (a CheckSSHResourceModelAdapter) ToAPIArgument(model CheckSSHResourceModel) (*upapi.CheckSSH, error) {
	api := upapi.CheckSSH{
		Name:                   model.Name.ValueString(),
		Address:                model.Address.ValueString(),
		Port:                   model.Port.ValueInt64(),
		Sensitivity:            model.Sensitivity.ValueInt64(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		NumRetries:             model.NumRetries.ValueInt64(),
		UseIpVersion:           model.UseIpVersion.ValueString(),
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

func (a CheckSSHResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckSSHResourceModel, error) {
	model := CheckSSHResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		Address:                types.StringValue(api.Address),
		Port:                   types.Int64Value(api.Port),
		Sensitivity:            types.Int64Value(api.Sensitivity),
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
	}
	return &model, nil
}

type CheckSSHResourceAPI struct {
	provider *providerImpl
}

func (m *CheckSSHResourceModel) APIResource() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

func (c CheckSSHResourceAPI) Create(ctx context.Context, arg upapi.CheckSSH) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateSSH(ctx, arg)
}

func (c CheckSSHResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckSSHResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckSSH) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateSSH(ctx, pk, arg)
}

func (c CheckSSHResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
