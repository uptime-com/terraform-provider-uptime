package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckCloudStatusResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[CheckCloudStatusResourceModel, upapi.CheckCloudStatus, upapi.Check](
		CheckCloudStatusResourceAPI{provider: p},
		CheckCloudStatusResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_cloudstatus",
			Schema: schema.Schema{
				Description: "Monitor a public cloud provider status feed (Cloud Status check). " +
					"Configure either a single legacy `service_name`, or a `group` plus `monitoring_type` " +
					"(`ALL` to track every service in the group, `SPECIFIC` to track entries listed in " +
					"`services` and/or `service_titles`). Use the upstream " +
					"`/api/v1/checks/cloudstatus-groups/` and `/api/v1/checks/cloudstatus-services/` " +
					"endpoints to discover valid IDs. " +
					"Import using the check ID: `terraform import uptime_check_cloudstatus.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"locations":      LocationsOptionalSchemaAttribute(p),
					"tags":           TagsSchemaAttribute(),
					"is_paused":      IsPausedSchemaAttribute(),
					"notify_only_on_down": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Opt out of maintenance notifications.",
					},
					"service_name": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
						Description: "Deprecated: legacy single-component identifier. Prefer `group` + " +
							"`monitoring_type`.",
					},
					"group": schema.Int64Attribute{
						Optional:    true,
						Description: "Cloud status group ID to monitor. Write-only on the server.",
					},
					"monitoring_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
						Description: "Selects how `group` is monitored: `ALL` for every service in the group, " +
							"`SPECIFIC` for entries listed in `services`/`service_titles`.",
					},
					"services": schema.SetAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.Int64Type,
						Description: "Specific service IDs to monitor when `monitoring_type` is `SPECIFIC`.",
					},
					"service_titles": schema.SetAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Description: "Service title strings; matching current and future services are " +
							"auto-monitored when `monitoring_type` is `SPECIFIC`.",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type CheckCloudStatusResourceModel struct {
	ID               types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL              types.String `tfsdk:"url" ref:"URL,opt"`
	Name             types.String `tfsdk:"name"`
	ContactGroups    types.Set    `tfsdk:"contact_groups"`
	Locations        types.Set    `tfsdk:"locations"`
	Tags             types.Set    `tfsdk:"tags"`
	IsPaused         types.Bool   `tfsdk:"is_paused"`
	NotifyOnlyOnDown types.Bool   `tfsdk:"notify_only_on_down"`
	ServiceName      types.String `tfsdk:"service_name"`
	Group            types.Int64  `tfsdk:"group"`
	MonitoringType   types.String `tfsdk:"monitoring_type"`
	Services         types.Set    `tfsdk:"services"`
	ServiceTitles    types.Set    `tfsdk:"service_titles"`
}

func (m CheckCloudStatusResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckCloudStatusResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
}

func (a CheckCloudStatusResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckCloudStatusResourceModel, diag.Diagnostics) {
	model := *new(CheckCloudStatusResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckCloudStatusResourceModelAdapter) ToAPIArgument(model CheckCloudStatusResourceModel) (*upapi.CheckCloudStatus, error) {
	cfg := upapi.CheckCloudStatusConfig{
		NotifyOnlyOnDown: model.NotifyOnlyOnDown.ValueBool(),
		ServiceName:      model.ServiceName.ValueString(),
		MonitoringType:   model.MonitoringType.ValueString(),
		Services:         setToInt64Slice(model.Services),
		ServiceTitles:    setToStringSlice(model.ServiceTitles),
	}
	if !model.Group.IsNull() && !model.Group.IsUnknown() {
		v := model.Group.ValueInt64()
		cfg.Group = &v
	}
	return &upapi.CheckCloudStatus{
		Name:              model.Name.ValueString(),
		ContactGroups:     a.ContactGroups(model.ContactGroups),
		Locations:         a.Locations(model.Locations),
		Tags:              a.Tags(model.Tags),
		IsPaused:          upapi.BoolPtr(model.IsPaused.ValueBool()),
		CloudStatusConfig: cfg,
	}, nil
}

func (a CheckCloudStatusResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckCloudStatusResourceModel, error) {
	model := CheckCloudStatusResourceModel{
		ID:               types.Int64Value(api.PK),
		URL:              types.StringValue(api.URL),
		Name:             types.StringValue(api.Name),
		ContactGroups:    a.ContactGroupsValue(api.ContactGroups),
		Locations:        a.LocationsValue(api.Locations),
		Tags:             a.TagsValue(api.Tags),
		IsPaused:         types.BoolValue(api.IsPaused),
		NotifyOnlyOnDown: types.BoolValue(false),
		ServiceName:      types.StringValue(""),
		// `group` is write-only server-side; preserve plan value upstream of state by
		// returning a null so Terraform keeps whatever the user set.
		Group:          types.Int64Null(),
		MonitoringType: types.StringValue(""),
		Services:       types.SetValueMust(types.Int64Type, []attr.Value{}),
		ServiceTitles:  types.SetValueMust(types.StringType, []attr.Value{}),
	}
	if api.CloudStatusConfig != nil {
		c := api.CloudStatusConfig
		model.NotifyOnlyOnDown = types.BoolValue(c.NotifyOnlyOnDown)
		model.ServiceName = types.StringValue(c.ServiceName)
		model.MonitoringType = types.StringValue(c.MonitoringType)
		model.Services = int64SliceToSet(c.Services)
		model.ServiceTitles = stringSliceToSet(c.ServiceTitles)
	}
	return &model, nil
}

// PreservePlanValues keeps the `group` write-only value from the plan since
// the API never echoes it back. Implements PlanValuePreserver.
func (a CheckCloudStatusResourceModelAdapter) PreservePlanValues(result *CheckCloudStatusResourceModel, plan *CheckCloudStatusResourceModel) *CheckCloudStatusResourceModel {
	if result == nil || plan == nil {
		return result
	}
	result.Group = plan.Group
	return result
}

type CheckCloudStatusResourceAPI struct {
	provider *providerImpl
}

func (c CheckCloudStatusResourceAPI) Create(ctx context.Context, arg upapi.CheckCloudStatus) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateCloudStatus(ctx, arg)
}

func (c CheckCloudStatusResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckCloudStatusResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckCloudStatus) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateCloudStatus(ctx, pk, arg)
}

func (c CheckCloudStatusResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}

func setToInt64Slice(s types.Set) []int64 {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	elems := s.Elements()
	out := make([]int64, 0, len(elems))
	for _, e := range elems {
		v, ok := e.(basetypes.Int64Value)
		if !ok {
			continue
		}
		out = append(out, v.ValueInt64())
	}
	return out
}

func setToStringSlice(s types.Set) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	elems := s.Elements()
	out := make([]string, 0, len(elems))
	for _, e := range elems {
		v, ok := e.(basetypes.StringValue)
		if !ok {
			continue
		}
		out = append(out, v.ValueString())
	}
	return out
}

func int64SliceToSet(in []int64) types.Set {
	vals := make([]attr.Value, 0, len(in))
	for _, v := range in {
		vals = append(vals, types.Int64Value(v))
	}
	return types.SetValueMust(types.Int64Type, vals)
}

func stringSliceToSet(in []string) types.Set {
	vals := make([]attr.Value, 0, len(in))
	for _, v := range in {
		vals = append(vals, types.StringValue(v))
	}
	return types.SetValueMust(types.StringType, vals)
}
