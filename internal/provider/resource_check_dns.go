package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shopspring/decimal"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckDNSResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckDNSResourceModel, upapi.CheckDNS, upapi.Check]{
		CheckDNSResourceAPI{provider: p},
		CheckDNSResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_dns",
			Schema: schema.Schema{
				Description: "Monitor for DNS failures or changes",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"locations":      LocationsSchemaAttribute(p),
					"tags":           TagsSchemaAttribute(),
					"is_paused":      IsPausedSchemaAttribute(),
					"interval":       IntervalSchemaAttribute(5),
					"threshold":      ThresholdSchemaAttribute(20),
					"address":        AddressHostnameSchemaAttribute(),
					"dns_server": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "DNS server to query",
					},
					"dns_record_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("ANY"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"ANY", "A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT"}),
						},
					},
					"expect_string": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "IP Address, Domain Name or String to expect in response",
					},
					"sensitivity":               SensitivitySchemaAttribute(2),
					"num_retries":               NumRetriesAttribute(2),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
				},
			},
		},
	}
}

type CheckDNSResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Locations              types.Set    `tfsdk:"locations"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	Threshold              types.Int64  `tfsdk:"threshold"`
	Address                types.String `tfsdk:"address"`
	DNSServer              types.String `tfsdk:"dns_server"`
	DNSRecordType          types.String `tfsdk:"dns_record_type"`
	ExpectString           types.String `tfsdk:"expect_string"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLA                    types.Object `tfsdk:"sla"`

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckDNSResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckDNSResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckDNSResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckDNSResourceModel, diag.Diagnostics) {
	model := *new(CheckDNSResourceModel)
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

func (a CheckDNSResourceModelAdapter) ToAPIArgument(model CheckDNSResourceModel) (_ *upapi.CheckDNS, err error) {
	obj := upapi.CheckDNS{
		Name:                   model.Name.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		Threshold:              model.Threshold.ValueInt64(),
		Address:                model.Address.ValueString(),
		DNSServer:              model.DNSServer.ValueString(),
		DNSRecordType:          model.DNSRecordType.ValueString(),
		ExpectString:           model.ExpectString.ValueString(),
		Sensitivity:            model.Sensitivity.ValueInt64(),
		NumRetries:             model.NumRetries.ValueInt64(),
		Notes:                  model.Notes.ValueString(),
		IncludeInGlobalMetrics: model.IncludeInGlobalMetrics.ValueBool(),
	}

	if model.sla != nil {
		if !model.sla.Uptime.IsUnknown() {
			obj.UptimeSLA = model.sla.Uptime.ValueDecimal()
		}
		if !model.sla.Latency.IsUnknown() {
			obj.ResponseTimeSLA = decimal.NewFromFloat(model.sla.Latency.ValueDuration().Seconds())
		}
	}

	return &obj, nil
}

func (a CheckDNSResourceModelAdapter) FromAPIResult(api upapi.Check) (_ *CheckDNSResourceModel, err error) {
	model := CheckDNSResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Locations:              a.LocationsValue(api.Locations),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		Threshold:              types.Int64Value(api.Threshold),
		Address:                types.StringValue(api.Address),
		DNSServer:              types.StringValue(api.DNSServer),
		DNSRecordType:          types.StringValue(api.DNSRecordType),
		ExpectString:           types.StringValue(api.ExpectString),
		Sensitivity:            types.Int64Value(api.Sensitivity),
		NumRetries:             types.Int64Value(api.NumRetries),
		Notes:                  types.StringValue(api.Notes),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Uptime:  DecimalValue(api.UptimeSLA),
			Latency: DurationValueFromDecimalSeconds(api.ResponseTimeSLA),
		}),
	}
	return &model, nil
}

type CheckDNSResourceAPI struct {
	provider *providerImpl
}

func (a CheckDNSResourceAPI) Create(ctx context.Context, arg upapi.CheckDNS) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateDNS(ctx, arg)
}

func (a CheckDNSResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a CheckDNSResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckDNS) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateDNS(ctx, pk, arg)
}

func (a CheckDNSResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
