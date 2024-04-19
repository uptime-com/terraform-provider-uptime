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

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckPageSpeedResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckPageSpeedResourceModel, upapi.CheckPageSpeed, upapi.Check]{
		api: CheckPageSpeedResourceAPI{provider: p},
		mod: CheckPageSpeedResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_pagespeed",
			Schema: schema.Schema{
				Description: "Page Speed Check",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"locations":      LocationsSchemaAttributeWithDefaults(p, "Dedicated-Canada-Toronto"),
					"tags":           TagsSchemaAttribute(),
					"is_paused":      IsPausedSchemaAttribute(),
					"interval":       IntervalSchemaAttribute(1440),
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
					"headers": schema.StringAttribute{
						Optional:  true,
						Computed:  true,
						Sensitive: true,
						Default:   stringdefault.StaticString(""),
					},
					"script":      ScriptSchemaAttribute(),
					"num_retries": NumRetriesSchemaAttribute(2),
					"notes":       NotesSchemaAttribute(),
					"config": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"emulated_device": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("DEFAULT"),
							},
							"connection_throttling": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("UNTHROTTLED"),
							},
							"exclude_urls": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"uptime_grade_threshold": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
						},
					},
				},
			},
		},
	}
}

type CheckPageSpeedResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	Locations     types.Set    `tfsdk:"locations"`
	Tags          types.Set    `tfsdk:"tags"`
	IsPaused      types.Bool   `tfsdk:"is_paused"`
	Interval      types.Int64  `tfsdk:"interval"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	Headers       types.String `tfsdk:"headers"`
	Script        RawJson      `tfsdk:"script"`
	NumRetries    types.Int64  `tfsdk:"num_retries"`
	Notes         types.String `tfsdk:"notes"`
	Config        types.Object `tfsdk:"config"`

	config *CheckPageSpeedConfigAttribute `tfsdk:"-"`
}

func (m CheckPageSpeedResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckPageSpeedConfigAttribute struct {
	EmulatedDevice       types.String `tfsdk:"emulated_device"`
	ConnectionThrottling types.String `tfsdk:"connection_throttling"`
	ExcludeURLs          types.String `tfsdk:"exclude_urls"`
	UptimeGradeThreshold types.String `tfsdk:"uptime_grade_threshold"`
}

type CheckPageSpeedResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
}

func (a CheckPageSpeedResourceModelAdapter) ConfigAttributeContext(ctx context.Context, v types.Object) (*CheckPageSpeedConfigAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := CheckPageSpeedConfigAttribute{}
	diags := v.As(ctx, &m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (a CheckPageSpeedResourceModelAdapter) configAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"emulated_device":        types.StringType,
		"connection_throttling":  types.StringType,
		"exclude_urls":           types.StringType,
		"uptime_grade_threshold": types.StringType,
	}
}

func (a CheckPageSpeedResourceModelAdapter) configAttributeValues(model CheckPageSpeedConfigAttribute) map[string]attr.Value {
	return map[string]attr.Value{
		"emulated_device":        model.EmulatedDevice,
		"connection_throttling":  model.ConnectionThrottling,
		"exclude_urls":           model.ExcludeURLs,
		"uptime_grade_threshold": model.UptimeGradeThreshold,
	}
}

func (a CheckPageSpeedResourceModelAdapter) ConfigAttributeValue(m CheckPageSpeedConfigAttribute) types.Object {
	return types.ObjectValueMust(a.configAttributeTypes(), a.configAttributeValues(m))
}

func (a CheckPageSpeedResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckPageSpeedResourceModel, diag.Diagnostics) {
	model := CheckPageSpeedResourceModel{}
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}

	model.config, diags = a.ConfigAttributeContext(ctx, model.Config)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckPageSpeedResourceModelAdapter) ToAPIArgument(model CheckPageSpeedResourceModel) (*upapi.CheckPageSpeed, error) {
	api := upapi.CheckPageSpeed{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroups(model.ContactGroups),
		Locations:     a.Locations(model.Locations),
		Tags:          a.Tags(model.Tags),
		IsPaused:      model.IsPaused.ValueBool(),
		Interval:      model.Interval.ValueInt64(),
		Username:      model.Username.ValueString(),
		Password:      model.Password.ValueString(),
		Headers:       model.Headers.ValueString(),
		Script:        model.Script.ValueString(),
		NumRetries:    model.NumRetries.ValueInt64(),
		Notes:         model.Notes.ValueString(),
	}
	if model.config != nil {
		api.Config = upapi.CheckPageSpeedConfig{
			EmulatedDevice:       model.config.EmulatedDevice.ValueString(),
			ConnectionThrottling: model.config.ConnectionThrottling.ValueString(),
			ExcludeURLs:          model.config.ExcludeURLs.ValueString(),
			UptimeGradeThreshold: model.config.UptimeGradeThreshold.ValueString(),
		}
	}
	return &api, nil
}

func (a CheckPageSpeedResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckPageSpeedResourceModel, error) {
	model := CheckPageSpeedResourceModel{
		ID:            types.Int64Value(api.PK),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsValue(api.ContactGroups),
		Locations:     a.LocationsValue(api.Locations),
		Tags:          a.TagsValue(api.Tags),
		IsPaused:      types.BoolValue(api.IsPaused),
		Interval:      types.Int64Value(api.Interval),
		Username:      types.StringValue(api.Username),
		Password:      types.StringValue(api.Password),
		Headers:       types.StringValue(api.Headers),
		Script:        RawJsonValue(api.Script),
		NumRetries:    types.Int64Value(api.NumRetries),
		Notes:         types.StringValue(api.Notes),
		Config: a.ConfigAttributeValue(CheckPageSpeedConfigAttribute{
			EmulatedDevice:       types.StringValue(api.PageSpeedConfig.EmulatedDevice),
			ConnectionThrottling: types.StringValue(api.PageSpeedConfig.ConnectionThrottling),
			ExcludeURLs:          types.StringValue(api.PageSpeedConfig.ExcludeURLs),
			UptimeGradeThreshold: types.StringValue(api.PageSpeedConfig.UptimeGradeThreshold),
		}),
	}
	return &model, nil
}

type CheckPageSpeedResourceAPI struct {
	provider *providerImpl
}

func (c CheckPageSpeedResourceAPI) Create(ctx context.Context, arg upapi.CheckPageSpeed) (*upapi.Check, error) {
	return c.provider.api.Checks().CreatePageSpeed(ctx, arg)
}

func (c CheckPageSpeedResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckPageSpeedResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckPageSpeed) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdatePageSpeed(ctx, pk, arg)
}

func (c CheckPageSpeedResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
