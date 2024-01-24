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

func NewCheckTransactionResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckTransactionResourceModel, upapi.CheckAPI, upapi.Check]{
		CheckTransactionResourceAPI{provider: p},
		CheckTransactionResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_transaction",
			Schema: schema.Schema{
				Description: "Transaction check to monitor your entire site by scanning for suitable checks to add",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"locations":                 LocationsSchemaAttribute(p),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"interval":                  IntervalSchemaAttribute(5),
					"threshold":                 ThresholdSchemaAttribute(30),
					"sensitivity":               SensitivitySchemaAttribute(2),
					"num_retries":               NumRetriesAttribute(2),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"script":                    ScriptSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
				},
			},
		},
	}
}

type CheckTransactionResource struct {
	*APIResource[CheckTransactionResourceModel, upapi.CheckTransaction, upapi.Check]
}

type CheckTransactionResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	URL                    types.String `tfsdk:"url"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Locations              types.Set    `tfsdk:"locations"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	Threshold              types.Int64  `tfsdk:"threshold"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`

	Script RawJson `tfsdk:"script"`

	SLA types.Object  `tfsdk:"sla"`
	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckTransactionResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckTransactionResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter

	SLAAttributeContextAdapter
}

func (a CheckTransactionResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckTransactionResourceModel, diag.Diagnostics) {
	model := *new(CheckTransactionResourceModel)
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

func (a CheckTransactionResourceModelAdapter) ToAPIArgument(model CheckTransactionResourceModel) (_ *upapi.CheckAPI, err error) {
	api := upapi.CheckAPI{
		Name:                   model.Name.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Locations:              a.Locations(model.Locations),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Interval:               model.Interval.ValueInt64(),
		Threshold:              model.Threshold.ValueInt64(),
		Script:                 model.Script.ValueString(),
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

func (a CheckTransactionResourceModelAdapter) FromAPIResult(api upapi.Check) (_ *CheckTransactionResourceModel, err error) {
	model := CheckTransactionResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Locations:              a.LocationsValue(api.Locations),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Interval:               types.Int64Value(api.Interval),
		Threshold:              types.Int64Value(api.Threshold),
		Script:                 RawJsonValue(api.Script),
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

var _ API[upapi.CheckAPI, upapi.Check] = (*CheckTransactionResourceAPI)(nil)

type CheckTransactionResourceAPI struct {
	provider *providerImpl
}

func (a CheckTransactionResourceAPI) Create(ctx context.Context, arg upapi.CheckAPI) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateAPI(ctx, arg)
}

func (a CheckTransactionResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a CheckTransactionResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckAPI) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateAPI(ctx, pk, arg)
}

func (a CheckTransactionResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
