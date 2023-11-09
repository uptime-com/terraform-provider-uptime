package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shopspring/decimal"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckHTTPResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckHTTPResourceModel, upapi.CheckHTTP, upapi.Check]{
		api: CheckHTTPResourceAPI{provider: p},
		mod: CheckHTTPResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_http",
			Schema: schema.Schema{
				Description: "Monitor a URL for specific status code(s)",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"address":                   AddressURLSchemaAttribute(),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"locations":                 LocationsSchemaAttribute(p.getLocations),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"interval":                  IntervalSchemaAttribute(5),
					"threshold":                 ThresholdSchemaAttribute(40),
					"sensitivity":               SensitivitySchemaAttribute(2),
					"num_retries":               NumRetriesAttribute(2),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
					"port": schema.Int64Attribute{
						Computed: true,
						Optional: true,
						Default:  int64default.StaticInt64(0),
					},
					"username": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Computed:  true,
						Sensitive: true,
						Default:   stringdefault.StaticString(""),
					},
					"proxy": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"status_code": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("200"),
					},
					"send_string": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "String to post",
					},
					"expect_string": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"expect_string_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("STRING"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"STRING", "REGEX", "INVERSE_REGEX"}),
						},
					},
					"headers": HeadersSchemaAttribute(),
					"version": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(2),
						Description: "Check version to use. Keep default value unless you are absolutely sure you need to change it",
					},
					"encryption": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("SSL_TLS"),
						Description: "Whether to verify SSL/TLS certificates",
					},
				},
			},
		},
	}
}

type CheckHTTPResourceModel struct {
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
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	Proxy                  types.String `tfsdk:"proxy"`
	StatusCode             types.String `tfsdk:"status_code"`
	SendString             types.String `tfsdk:"send_string"`
	ExpectString           types.String `tfsdk:"expect_string"`
	ExpectStringType       types.String `tfsdk:"expect_string_type"`
	Encryption             types.String `tfsdk:"encryption"`
	Threshold              types.Int64  `tfsdk:"threshold"`
	Headers                types.Map    `tfsdk:"headers"`
	Version                types.Int64  `tfsdk:"version"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLA                    types.Object `tfsdk:"sla"`

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckHTTPResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckHTTPResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckHTTPResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckHTTPResourceModel, diag.Diagnostics) {
	model := *new(CheckHTTPResourceModel)
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

func (a CheckHTTPResourceModelAdapter) ToAPIArgument(model CheckHTTPResourceModel) (_ *upapi.CheckHTTP, err error) {
	api := upapi.CheckHTTP{
		Name:                   model.Name.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		Address:                model.Address.ValueString(),
		Port:                   model.Port.ValueInt64(),
		Username:               model.Username.ValueString(),
		Password:               model.Password.ValueString(),
		Proxy:                  model.Proxy.ValueString(),
		StatusCode:             model.StatusCode.ValueString(),
		SendString:             model.SendString.ValueString(),
		ExpectString:           model.ExpectString.ValueString(),
		ExpectStringType:       model.ExpectStringType.ValueString(),
		Encryption:             model.Encryption.ValueString(),
		Threshold:              model.Threshold.ValueInt64(),
		Version:                model.Version.ValueInt64(),
		Sensitivity:            model.Sensitivity.ValueInt64(),
		NumRetries:             model.NumRetries.ValueInt64(),
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

func (a CheckHTTPResourceModelAdapter) FromAPIResult(api upapi.Check) (_ *CheckHTTPResourceModel, err error) {
	model := CheckHTTPResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Locations:              a.LocationsValue(api.Locations),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		Address:                types.StringValue(api.Address),
		Port:                   types.Int64Value(api.Port),
		Username:               types.StringValue(api.Username),
		Password:               types.StringValue(api.Password),
		Headers:                types.MapValueMust(types.ListType{ElemType: types.StringType}, nil), // TODO: fix this
		Proxy:                  types.StringValue(api.Proxy),
		StatusCode:             types.StringValue(api.StatusCode),
		SendString:             types.StringValue(api.SendString),
		ExpectString:           types.StringValue(api.ExpectString),
		ExpectStringType:       types.StringValue(api.ExpectStringType),
		Encryption:             types.StringValue(api.Encryption),
		Threshold:              types.Int64Value(api.Threshold),
		Version:                types.Int64Value(api.Version),
		Sensitivity:            types.Int64Value(api.Sensitivity),
		NumRetries:             types.Int64Value(api.NumRetries),
		Notes:                  types.StringValue(api.Notes),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Latency: DurationValueFromDecimalSeconds(api.ResponseTimeSLA),
			Uptime:  DecimalValue(api.UptimeSLA),
		}),
	}
	return &model, nil
}

type CheckHTTPResourceAPI struct {
	provider *providerImpl
}

func (a CheckHTTPResourceAPI) Create(ctx context.Context, arg upapi.CheckHTTP) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateHTTP(ctx, arg)
}

func (a CheckHTTPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a CheckHTTPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckHTTP) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateHTTP(ctx, pk, arg)
}

func (a CheckHTTPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
