package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewDashboardResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[DashboardResourceModel, upapi.Dashboard, upapi.Dashboard]{
		api: DashboardResourceAPI{provider: p},
		mod: DashboardResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "dashboard",
			Schema: schema.Schema{
				Description: "Custom dashboard resource",
				Attributes: map[string]schema.Attribute{
					"id":   IDSchemaAttribute(),
					"name": NameSchemaAttribute(),
					"ordering": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Where to place the dashboard in the list",
						Default:     int64default.StaticInt64(10),
					},
					"is_pinned": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Whether the dashboard is pinned to the top of the dashboard list",
					},
					"metrics": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Metrics related attributes",
						Attributes: map[string]schema.Attribute{
							"show_section": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Whether to show the metrics section",
							},
							"for_all_checks": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Whether to show metrics for all checks",
							},
						},
					},
					"services": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Services related attributes",
						Attributes: map[string]schema.Attribute{
							"show_section": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Whether to show the services section",
							},
							"num_to_show": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "The number of services to show",
								Default:     int64default.StaticInt64(4),
							},
							"include": schema.SingleNestedAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Which services to include",
								Attributes: map[string]schema.Attribute{
									"up": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to include up services",
									},
									"down": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to include down services",
									},
									"paused": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to include paused services",
									},
									"maintenance": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to include services in maintenance mode",
									},
								},
							},
							"sort": schema.SingleNestedAttribute{
								Optional:    true,
								Computed:    true,
								Description: "How to sort services",
								Attributes: map[string]schema.Attribute{
									"primary": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The primary sort for services",
										Default:     stringdefault.StaticString("is_paused,cached_state_is_up"),
									},
									"secondary": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The secondary sort for services",
										Default:     stringdefault.StaticString("-cached_last_down_alert_at"),
									},
								},
							},
							"show": schema.SingleNestedAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Which service attributes to show",
								Attributes: map[string]schema.Attribute{
									"uptime": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to show uptime for services",
										Default:     booldefault.StaticBool(true),
									},
									"response_time": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to show response time for services",
										Default:     booldefault.StaticBool(true),
									},
								},
							},
						},
					},
					"alerts": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Alerts related attributes",
						Attributes: map[string]schema.Attribute{
							"show_section": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Whether to show the alerts section",
								Default:     booldefault.StaticBool(true),
							},
							"num_to_show": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "The number of alerts to show",
								Default:     int64default.StaticInt64(10),
							},
							"for_all_checks": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Whether to show alerts for all checks",
							},
							"include": schema.SingleNestedAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Alerts related attributes",
								Attributes: map[string]schema.Attribute{
									"ignored": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to include ignored alerts",
									},
									"resolved": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Whether to include resolved alerts",
									},
								},
							},
						},
					},
					"selected": schema.SingleNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Selected services to show on the dashboard",
						Attributes: map[string]schema.Attribute{
							"services": schema.SetAttribute{
								Description: "The services collection to show on the dashboard",
								ElementType: types.StringType,
								Computed:    true,
								Optional:    true,
								Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
							},
							"tags": TagsSchemaAttribute(),
						},
					},
				},
			},
		},
	}
}

type DashboardResourceModel struct {
	ID       types.Int64  `tfsdk:"id" ref:"PK,opt"`
	Name     types.String `tfsdk:"name"`
	Ordering types.Int64  `tfsdk:"ordering"`
	IsPinned types.Bool   `tfsdk:"is_pinned"`
	Metrics  types.Object `tfsdk:"metrics"`
	Services types.Object `tfsdk:"services"`
	Alerts   types.Object `tfsdk:"alerts"`
	Selected types.Object `tfsdk:"selected"`

	metrics  *DashboardMetricsAttribute  `tfsdk:"-"`
	services *DashboardServicesAttribute `tfsdk:"-"`
	alerts   *DashboardAlertsAttribute   `tfsdk:"-"`
	selected *DashboardSelectedAttribute `tfsdk:"-"`
}

func (m DashboardResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type DashboardMetricsAttribute struct {
	ShowSection  types.Bool `tfsdk:"show_section"`
	ForAllChecks types.Bool `tfsdk:"for_all_checks"`
}

type DashboardServicesAttribute struct {
	ShowSection types.Bool   `tfsdk:"show_section"`
	NumToShow   types.Int64  `tfsdk:"num_to_show"`
	Include     types.Object `tfsdk:"include"`
	Sort        types.Object `tfsdk:"sort"`
	Show        types.Object `tfsdk:"show"`

	include *DashboardServicesIncludeAttribute `tfsdk:"-"`
	sort    *DashboardServicesSortAttribute    `tfsdk:"-"`
	show    *DashboardServicesShowAttribute    `tfsdk:"-"`
}

type DashboardServicesIncludeAttribute struct {
	Up          types.Bool `tfsdk:"up"`
	Down        types.Bool `tfsdk:"down"`
	Paused      types.Bool `tfsdk:"paused"`
	Maintenance types.Bool `tfsdk:"maintenance"`
}

type DashboardServicesSortAttribute struct {
	Primary   types.String `tfsdk:"primary"`
	Secondary types.String `tfsdk:"secondary"`
}

type DashboardServicesShowAttribute struct {
	Uptime       types.Bool `tfsdk:"uptime"`
	ResponseTime types.Bool `tfsdk:"response_time"`
}

type DashboardAlertsAttribute struct {
	ShowSection  types.Bool   `tfsdk:"show_section"`
	NumToShow    types.Int64  `tfsdk:"num_to_show"`
	ForAllChecks types.Bool   `tfsdk:"for_all_checks"`
	Include      types.Object `tfsdk:"include"`

	include *DashboardAlertsIncludeAttribute `tfsdk:"-"`
}

type DashboardAlertsIncludeAttribute struct {
	Ignored  types.Bool `tfsdk:"ignored"`
	Resolved types.Bool `tfsdk:"resolved"`
}

type DashboardSelectedAttribute struct {
	Services types.Set `tfsdk:"services"`
	Tags     types.Set `tfsdk:"tags"`
}

type DashboardResourceModelAdapter struct {
	SetAttributeAdapter[string]
}

func (a DashboardResourceModelAdapter) metricsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"show_section":   types.BoolType,
		"for_all_checks": types.BoolType,
	}
}

func (a DashboardResourceModelAdapter) metricsAttributeValues(model DashboardMetricsAttribute) map[string]attr.Value {
	v := map[string]attr.Value{
		"show_section":   model.ShowSection,
		"for_all_checks": model.ForAllChecks,
	}
	return v
}

func (a DashboardResourceModelAdapter) MetricsAttributeContext(ctx context.Context, v types.Object) (*DashboardMetricsAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := new(DashboardMetricsAttribute)
	diags := v.As(ctx, m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return m, nil
}

func (a DashboardResourceModelAdapter) MetricsAttributeValue(model DashboardMetricsAttribute) types.Object {
	return types.ObjectValueMust(a.metricsAttributeTypes(), a.metricsAttributeValues(model))
}

func (a DashboardResourceModelAdapter) servicesAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"show_section": types.BoolType,
		"num_to_show":  types.Int64Type,
		"include":      types.ObjectType{}.WithAttributeTypes(a.servicesIncludeAttributeTypes()),
		"sort":         types.ObjectType{}.WithAttributeTypes(a.servicesSortAttributeTypes()),
		"show":         types.ObjectType{}.WithAttributeTypes(a.servicesShowAttributeTypes()),
	}
}

func (a DashboardResourceModelAdapter) servicesAttributeValues(model DashboardServicesAttribute) map[string]attr.Value {
	v := map[string]attr.Value{
		"show_section": model.ShowSection,
		"num_to_show":  model.NumToShow,
		"include":      model.Include,
		"show":         model.Show,
		"sort":         model.Sort,
	}
	return v
}

func (a DashboardResourceModelAdapter) ServicesAttributeContext(ctx context.Context, v types.Object) (*DashboardServicesAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := *new(DashboardServicesAttribute)
	diags := v.As(ctx, &m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	m.show, diags = a.ServicesShowAttributeContext(ctx, m.Show)
	if diags.HasError() {
		return nil, diags
	}
	m.sort, diags = a.ServicesSortAttributeContext(ctx, m.Sort)
	if diags.HasError() {
		return nil, diags
	}
	m.include, diags = a.ServicesIncludeAttributeContext(ctx, m.Include)
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (a DashboardResourceModelAdapter) ServicesAttributeValue(model DashboardServicesAttribute) types.Object {
	return types.ObjectValueMust(a.servicesAttributeTypes(), a.servicesAttributeValues(model))
}

func (a DashboardResourceModelAdapter) servicesSortAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"primary":   types.StringType,
		"secondary": types.StringType,
	}
}

func (a DashboardResourceModelAdapter) servicesSortAttributeValues(
	model DashboardServicesSortAttribute,
) map[string]attr.Value {
	return map[string]attr.Value{
		"primary":   model.Primary,
		"secondary": model.Secondary,
	}
}

func (a DashboardResourceModelAdapter) ServicesSortAttributeContext(ctx context.Context, v types.Object) (*DashboardServicesSortAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := new(DashboardServicesSortAttribute)
	diags := v.As(ctx, m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return m, nil
}

func (a DashboardResourceModelAdapter) ServicesSortAttributeValue(model DashboardServicesSortAttribute) types.Object {
	return types.ObjectValueMust(a.servicesSortAttributeTypes(), a.servicesSortAttributeValues(model))
}

func (a DashboardResourceModelAdapter) servicesShowAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"uptime":        types.BoolType,
		"response_time": types.BoolType,
	}
}

func (a DashboardResourceModelAdapter) servicesShowAttributeValues(
	model DashboardServicesShowAttribute,
) map[string]attr.Value {
	return map[string]attr.Value{
		"uptime":        model.Uptime,
		"response_time": model.ResponseTime,
	}
}

func (a DashboardResourceModelAdapter) ServicesShowAttributeContext(ctx context.Context, v types.Object) (*DashboardServicesShowAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := *new(DashboardServicesShowAttribute)
	diags := v.As(ctx, &m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (a DashboardResourceModelAdapter) ServicesShowAttributeValue(model DashboardServicesShowAttribute) types.Object {
	return types.ObjectValueMust(a.servicesShowAttributeTypes(), a.servicesShowAttributeValues(model))
}

func (a DashboardResourceModelAdapter) servicesIncludeAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"up":          types.BoolType,
		"down":        types.BoolType,
		"paused":      types.BoolType,
		"maintenance": types.BoolType,
	}
}

func (a DashboardResourceModelAdapter) servicesIncludeAttributeValues(
	model DashboardServicesIncludeAttribute,
) map[string]attr.Value {
	return map[string]attr.Value{
		"up":          model.Up,
		"down":        model.Down,
		"paused":      model.Paused,
		"maintenance": model.Maintenance,
	}
}

func (a DashboardResourceModelAdapter) ServicesIncludeAttributeContext(ctx context.Context, v types.Object) (*DashboardServicesIncludeAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := *new(DashboardServicesIncludeAttribute)
	diags := v.As(ctx, &m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (a DashboardResourceModelAdapter) ServicesIncludeAttributeValue(model DashboardServicesIncludeAttribute) types.Object {
	return types.ObjectValueMust(a.servicesIncludeAttributeTypes(), a.servicesIncludeAttributeValues(model))
}

func (a DashboardResourceModelAdapter) alertsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"show_section":   types.BoolType,
		"num_to_show":    types.Int64Type,
		"for_all_checks": types.BoolType,
		"include":        types.ObjectType{}.WithAttributeTypes(a.alertsIncludeAttributeTypes()),
	}
}

func (a DashboardResourceModelAdapter) alertsAttributeValues(model DashboardAlertsAttribute) map[string]attr.Value {
	v := map[string]attr.Value{
		"show_section":   model.ShowSection,
		"num_to_show":    model.NumToShow,
		"for_all_checks": model.ForAllChecks,
		"include":        model.Include,
	}
	return v
}

func (a DashboardResourceModelAdapter) AlertsAttributeContext(ctx context.Context, v types.Object) (*DashboardAlertsAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := new(DashboardAlertsAttribute)
	diags := v.As(ctx, m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	m.include, diags = a.AlertsIncludeAttributeContext(ctx, m.Include)
	if diags.HasError() {
		return nil, diags
	}
	return m, nil
}

func (a DashboardResourceModelAdapter) AlertsAttributeValue(model DashboardAlertsAttribute) types.Object {
	return types.ObjectValueMust(a.alertsAttributeTypes(), a.alertsAttributeValues(model))
}

func (a DashboardResourceModelAdapter) alertsIncludeAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ignored":  types.BoolType,
		"resolved": types.BoolType,
	}
}

func (a DashboardResourceModelAdapter) alertsIncludeAttributeValues(
	model DashboardAlertsIncludeAttribute,
) map[string]attr.Value {
	return map[string]attr.Value{
		"ignored":  model.Ignored,
		"resolved": model.Resolved,
	}
}

func (a DashboardResourceModelAdapter) AlertsIncludeAttributeContext(ctx context.Context, v types.Object) (*DashboardAlertsIncludeAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := new(DashboardAlertsIncludeAttribute)
	diags := v.As(ctx, m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return m, nil
}

func (a DashboardResourceModelAdapter) AlertsIncludeAttributeValue(model DashboardAlertsIncludeAttribute) types.Object {
	return types.ObjectValueMust(a.alertsIncludeAttributeTypes(), a.alertsIncludeAttributeValues(model))
}

func (a DashboardResourceModelAdapter) selectedAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"services": types.SetType{}.WithElementType(types.StringType),
		"tags":     types.SetType{}.WithElementType(types.StringType),
	}
}

func (a DashboardResourceModelAdapter) selectedAttributeValues(model DashboardSelectedAttribute) map[string]attr.Value {
	v := map[string]attr.Value{
		"services": types.SetValueMust(types.StringType, model.Services.Elements()),
		"tags":     types.SetValueMust(types.StringType, model.Tags.Elements()),
	}
	return v
}

func (a DashboardResourceModelAdapter) SelectedAttributeContext(ctx context.Context, v types.Object) (*DashboardSelectedAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := *new(DashboardSelectedAttribute)
	diags := v.As(ctx, &m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (a DashboardResourceModelAdapter) SelectedAttributeValue(model DashboardSelectedAttribute) types.Object {
	return types.ObjectValueMust(a.selectedAttributeTypes(), a.selectedAttributeValues(model))
}

func (a DashboardResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*DashboardResourceModel, diag.Diagnostics) {
	model := *new(DashboardResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.metrics, diags = a.MetricsAttributeContext(ctx, model.Metrics)
	if diags.HasError() {
		return nil, diags
	}
	model.alerts, diags = a.AlertsAttributeContext(ctx, model.Alerts)
	if diags.HasError() {
		return nil, diags
	}
	model.services, diags = a.ServicesAttributeContext(ctx, model.Services)
	if diags.HasError() {
		return nil, diags
	}
	model.selected, diags = a.SelectedAttributeContext(ctx, model.Selected)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a DashboardResourceModelAdapter) ToAPIArgument(model DashboardResourceModel) (*upapi.Dashboard, error) {
	api := upapi.Dashboard{
		Name:     model.Name.ValueString(),
		IsPinned: model.IsPinned.ValueBool(),
		Ordering: model.Ordering.ValueInt64(),
	}
	if model.metrics != nil {
		api.MetricsShowSection = model.metrics.ShowSection.ValueBool()
		api.MetricsForAllChecks = model.metrics.ForAllChecks.ValueBool()
	}
	if model.services != nil {
		api.ServicesShowSection = model.services.ShowSection.ValueBool()
		api.ServicesNumToShow = model.services.NumToShow.ValueInt64()
		if model.services.include != nil {
			api.ServicesIncludeUp = model.services.include.Up.ValueBool()
			api.ServicesIncludeDown = model.services.include.Down.ValueBool()
			api.ServicesIncludePaused = model.services.include.Paused.ValueBool()
			api.ServicesIncludeMaintenance = model.services.include.Maintenance.ValueBool()
		}
		if model.services.sort != nil {
			api.ServicesPrimarySort = model.services.sort.Primary.ValueString()
			api.ServicesSecondarySort = model.services.sort.Secondary.ValueString()
		}
		if model.services.show != nil {
			api.ServicesShowUptime = model.services.show.Uptime.ValueBool()
			api.ServicesShowResponseTime = model.services.show.ResponseTime.ValueBool()
		}
	}
	if model.alerts != nil {
		api.AlertsnumToShow = model.alerts.NumToShow.ValueInt64()
		api.AlertsShowSection = model.alerts.ShowSection.ValueBool()
		api.AlertsForAllChecks = model.alerts.ForAllChecks.ValueBool()
		if model.alerts.include != nil {
			api.AlertsIncludeIgnored = model.alerts.include.Ignored.ValueBool()
			api.AlertsincludeResolved = model.alerts.include.Resolved.ValueBool()
		}
	}
	if model.selected != nil {
		api.ServicesSelected = a.Slice(model.selected.Services)
		api.ServicesTags = a.Slice(model.selected.Tags)
	}
	return &api, nil
}

func (a DashboardResourceModelAdapter) FromAPIResult(api upapi.Dashboard) (*DashboardResourceModel, error) {
	model := DashboardResourceModel{
		ID:       types.Int64Value(api.PK),
		Name:     types.StringValue(api.Name),
		Ordering: types.Int64Value(api.Ordering),
		IsPinned: types.BoolValue(api.IsPinned),
		Metrics: a.MetricsAttributeValue(DashboardMetricsAttribute{
			ShowSection:  types.BoolValue(api.MetricsShowSection),
			ForAllChecks: types.BoolValue(api.MetricsForAllChecks),
		}),
		Services: a.ServicesAttributeValue(DashboardServicesAttribute{
			ShowSection: types.BoolValue(api.ServicesShowSection),
			NumToShow:   types.Int64Value(api.ServicesNumToShow),
			Include: a.ServicesIncludeAttributeValue(DashboardServicesIncludeAttribute{
				Up:          types.BoolValue(api.ServicesIncludeUp),
				Down:        types.BoolValue(api.ServicesIncludeDown),
				Paused:      types.BoolValue(api.ServicesIncludePaused),
				Maintenance: types.BoolValue(api.ServicesIncludeMaintenance),
			}),
			Sort: a.ServicesSortAttributeValue(DashboardServicesSortAttribute{
				Primary:   types.StringValue(api.ServicesPrimarySort),
				Secondary: types.StringValue(api.ServicesSecondarySort),
			}),
			Show: a.ServicesShowAttributeValue(DashboardServicesShowAttribute{
				Uptime:       types.BoolValue(api.ServicesShowUptime),
				ResponseTime: types.BoolValue(api.ServicesShowResponseTime),
			}),
		}),
		Alerts: a.AlertsAttributeValue(DashboardAlertsAttribute{
			ShowSection:  types.BoolValue(api.AlertsShowSection),
			NumToShow:    types.Int64Value(api.AlertsnumToShow),
			ForAllChecks: types.BoolValue(api.AlertsForAllChecks),
			Include: a.AlertsIncludeAttributeValue(DashboardAlertsIncludeAttribute{
				Ignored:  types.BoolValue(api.AlertsIncludeIgnored),
				Resolved: types.BoolValue(api.AlertsincludeResolved),
			}),
		}),
	}

	var selectedAttribute DashboardSelectedAttribute
	if api.ServicesSelected != nil {
		selectedServices := []attr.Value{}
		for _, s := range api.ServicesSelected {
			selectedServices = append(selectedServices, types.StringValue(s))
		}
		selectedAttribute.Services = types.SetValueMust(types.StringType, selectedServices)
	}
	if api.ServicesTags != nil {
		selectedTags := []attr.Value{}
		for _, t := range api.ServicesTags {
			selectedTags = append(selectedTags, types.StringValue(t))
		}
		selectedAttribute.Tags = types.SetValueMust(types.StringType, selectedTags)
	}
	if len(selectedAttribute.Services.Elements()) > 0 || len(selectedAttribute.Tags.Elements()) > 0 {
		model.Selected = a.SelectedAttributeValue(selectedAttribute)
	}
	return &model, nil
}

type DashboardResourceAPI struct {
	provider *providerImpl
}

func (c DashboardResourceAPI) Create(ctx context.Context, arg upapi.Dashboard) (*upapi.Dashboard, error) {
	obj, err := c.provider.api.Dashboards().Create(ctx, arg)
	return obj, err
}

func (c DashboardResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Dashboard, error) {
	obj, err := c.provider.api.Dashboards().Get(ctx, pk)
	return obj, err
}

func (c DashboardResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.Dashboard) (*upapi.Dashboard, error) {
	obj, err := c.provider.api.Dashboards().Update(ctx, pk, arg)
	return obj, err
}

func (c DashboardResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Dashboards().Delete(ctx, pk)
}
