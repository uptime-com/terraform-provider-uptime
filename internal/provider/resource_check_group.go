package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/shopspring/decimal"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckGroupResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckGroupResourceModel, upapi.CheckGroup, upapi.Check]{
		CheckGroupResourceAPI{provider: p},
		CheckGroupResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_group",
			Schema: schema.Schema{
				Description: "Combine multiple checks",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"config": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"services": ServicesSchemaAttribute(),
							"tags":     TagsSchemaAttribute(),
							"down_condition": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("ANY"),
								Description: "Condition that determines when the group check is considered DOWN. " +
									"Valid values: `ANY`, `TWO`, `THREE`, `FOUR`, `FIVE`, `TEN`, " +
									"`ONE_PCT`, `THREE_PCT`, `FIVE_PCT`, `TEN_PCT`, `TWENTYFIVE_PCT`, `FIFTY_PCT`, `ALL`. " +
									"Numeric values (TWO-TEN) mean the group is DOWN when that many checks are down. " +
									"Percentage values (ONE_PCT-FIFTY_PCT) mean the group is DOWN when that percentage of checks are down. " +
									"Defaults to `ANY`.",
							},
							"uptime_percent_calculation": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("UP_DOWN_STATES"),
								Description: "Method used to calculate the group's uptime percentage. Valid values: `UP_DOWN_STATES` (calculates based on up/down state transitions), `AVERAGE` (calculates as the average uptime of all included checks). Defaults to `UP_DOWN_STATES`.",
							},
							"response_time": schema.SingleNestedAttribute{
								Optional: true,
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"calculation_mode": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Default:  stringdefault.StaticString("NONE"),
									},
									"check_type": schema.StringAttribute{
										Optional: true,
										Computed: true,
										Default:  stringdefault.StaticString("HTTP"),
									},
									"single_check": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type CheckGroupConfigTimeAttribute struct {
	CalculationMode types.String `tfsdk:"calculation_mode"`
	CheckType       types.String `tfsdk:"check_type"`
	SingleCheck     types.String `tfsdk:"single_check"`
}

type CheckGroupConfigAttribute struct {
	Services                 types.Set    `tfsdk:"services"`
	Tags                     types.Set    `tfsdk:"tags"`
	DownCondition            types.String `tfsdk:"down_condition"`
	UptimePercentCalculation types.String `tfsdk:"uptime_percent_calculation"`
	ResponseTime             types.Object `tfsdk:"response_time"`

	time *CheckGroupConfigTimeAttribute `tfsdk:"-"`
}

type CheckGroupResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLA                    types.Object `tfsdk:"sla"`
	Config                 types.Object `tfsdk:"config"`

	sla    *SLAAttribute              `tfsdk:"-"`
	config *CheckGroupConfigAttribute `tfsdk:"-"`
}

func (m CheckGroupResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

func (a CheckGroupResourceModelAdapter) configTimeAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"calculation_mode": types.StringType,
		"check_type":       types.StringType,
		"single_check":     types.StringType,
	}
}

func (a CheckGroupResourceModelAdapter) configTimeAttributeValues(m CheckGroupConfigTimeAttribute) map[string]attr.Value {
	return map[string]attr.Value{
		"calculation_mode": m.CalculationMode,
		"check_type":       m.CheckType,
		"single_check":     m.SingleCheck,
	}
}

func (a CheckGroupResourceModelAdapter) ConfigTimeAttributeContext(ctx context.Context, v types.Object) (*CheckGroupConfigTimeAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := new(CheckGroupConfigTimeAttribute)
	diags := v.As(ctx, m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return m, nil
}

func (a CheckGroupResourceModelAdapter) ConfigTimeAttributeValue(m CheckGroupConfigTimeAttribute) types.Object {
	return types.ObjectValueMust(a.configTimeAttributeTypes(), a.configTimeAttributeValues(m))
}

func (a CheckGroupResourceModelAdapter) groupConfigAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"services":                   types.SetType{}.WithElementType(types.StringType),
		"tags":                       types.SetType{}.WithElementType(types.StringType),
		"down_condition":             types.StringType,
		"uptime_percent_calculation": types.StringType,
		"response_time":              types.ObjectType{}.WithAttributeTypes(a.configTimeAttributeTypes()),
	}
}

func (a CheckGroupResourceModelAdapter) groupConfigAttributeValues(m CheckGroupConfigAttribute) map[string]attr.Value {
	return map[string]attr.Value{
		"services":                   types.SetValueMust(types.StringType, m.Services.Elements()),
		"tags":                       types.SetValueMust(types.StringType, m.Tags.Elements()),
		"down_condition":             m.DownCondition,
		"uptime_percent_calculation": m.UptimePercentCalculation,
		"response_time":              m.ResponseTime,
	}
}

func (a CheckGroupResourceModelAdapter) CheckGroupConfigAttributeValue(m CheckGroupConfigAttribute) types.Object {
	return types.ObjectValueMust(a.groupConfigAttributeTypes(), a.groupConfigAttributeValues(m))
}

func (a CheckGroupResourceModelAdapter) CheckGroupConfigAttributeContext(ctx context.Context, v types.Object) (*CheckGroupConfigAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := new(CheckGroupConfigAttribute)
	diags := v.As(ctx, m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	m.time, diags = a.ConfigTimeAttributeContext(ctx, m.ResponseTime)
	if diags.HasError() {
		return nil, diags
	}
	return m, nil
}

type CheckGroupResourceModelAdapter struct {
	SetAttributeAdapter[string]
	ContactGroupsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckGroupResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckGroupResourceModel, diag.Diagnostics) {
	model := *new(CheckGroupResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.sla, diags = a.SLAAttributeContext(ctx, model.SLA)
	if diags.HasError() {
		return nil, diags
	}
	model.config, diags = a.CheckGroupConfigAttributeContext(ctx, model.Config)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckGroupResourceModelAdapter) ToAPIArgument(model CheckGroupResourceModel) (_ *upapi.CheckGroup, err error) {
	obj := upapi.CheckGroup{
		Name:                   model.Name.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
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

	if model.config != nil {
		obj.Config.CheckServices = a.Slice(model.config.Services)
		obj.Config.CheckTags = a.Slice(model.config.Tags)
		obj.Config.UptimePercentCalculation = model.config.UptimePercentCalculation.ValueString()
		obj.Config.CheckDownCondition = model.config.DownCondition.ValueString()
		if model.config.time != nil {
			obj.Config.ResponseTimeCheckType = model.config.time.CheckType.ValueString()
			obj.Config.ResponseTimeSingleCheck = model.config.time.SingleCheck.ValueString()
			obj.Config.ResponseTimeCalculationMode = model.config.time.CalculationMode.ValueString()
		}
	}

	return &obj, nil
}

func (a CheckGroupResourceModelAdapter) FromAPIResult(api upapi.Check) (_ *CheckGroupResourceModel, err error) {
	model := CheckGroupResourceModel{
		ID:                     types.Int64Value(api.PK),
		Name:                   types.StringValue(api.Name),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Notes:                  types.StringValue(api.Notes),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Uptime:  DecimalValue(api.UptimeSLA),
			Latency: DurationValueFromDecimalSeconds(api.ResponseTimeSLA),
		}),
	}

	if api.GroupConfig != nil {
		configAttribute := CheckGroupConfigAttribute{
			DownCondition:            types.StringValue(api.GroupConfig.CheckDownCondition),
			UptimePercentCalculation: types.StringValue(api.GroupConfig.UptimePercentCalculation),
			ResponseTime: a.ConfigTimeAttributeValue(CheckGroupConfigTimeAttribute{
				CalculationMode: types.StringValue(api.GroupConfig.ResponseTimeCalculationMode),
				CheckType:       types.StringValue(api.GroupConfig.ResponseTimeCheckType),
				SingleCheck:     types.StringValue(api.GroupConfig.ResponseTimeSingleCheck),
			}),
		}
		if api.GroupConfig.CheckServices != nil {
			selectedServices := []attr.Value{}
			for _, service := range api.GroupConfig.CheckServices {
				selectedServices = append(selectedServices, types.StringValue(service))
			}
			configAttribute.Services = types.SetValueMust(types.StringType, selectedServices)
		}
		if api.GroupConfig.CheckTags != nil {
			selectedTags := []attr.Value{}
			for _, tag := range api.GroupConfig.CheckTags {
				selectedTags = append(selectedTags, types.StringValue(tag))
			}
			configAttribute.Tags = types.SetValueMust(types.StringType, selectedTags)
		}
		model.Config = a.CheckGroupConfigAttributeValue(configAttribute)
	}
	return &model, nil
}

// PreservePlanValues preserves the planned services values since the API
// returns check names instead of the numeric IDs that were sent.
// This implements the PlanValuePreserver interface.
func (a CheckGroupResourceModelAdapter) PreservePlanValues(result *CheckGroupResourceModel, plan *CheckGroupResourceModel) *CheckGroupResourceModel {
	if plan == nil || result == nil {
		return result
	}

	// If plan has config with services, preserve them in the result
	// because the API returns check names but we need to maintain the user's input (IDs)
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() && !result.Config.IsNull() {
		var planConfig, resultConfig CheckGroupConfigAttribute
		if diags := plan.Config.As(context.Background(), &planConfig, basetypes.ObjectAsOptions{}); diags.HasError() {
			return result
		}
		if diags := result.Config.As(context.Background(), &resultConfig, basetypes.ObjectAsOptions{}); diags.HasError() {
			return result
		}

		// Preserve planned services if they exist
		if !planConfig.Services.IsNull() && !planConfig.Services.IsUnknown() {
			resultConfig.Services = planConfig.Services
			result.Config = a.CheckGroupConfigAttributeValue(resultConfig)
		}
	}

	return result
}

type CheckGroupResourceAPI struct {
	provider *providerImpl
}

func (a CheckGroupResourceAPI) Create(ctx context.Context, arg upapi.CheckGroup) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateGroup(ctx, arg)
}

func (a CheckGroupResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a CheckGroupResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckGroup) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateGroup(ctx, pk, arg)
}

func (a CheckGroupResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
