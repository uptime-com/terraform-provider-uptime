package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewDashboardsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return DashboardsDataSource{p: p}
}

// DashboardsDataSchema defines the schema for the dashboards data source.
var DashboardsDataSchema = schema.Schema{
	Description: "Retrieve a list of all dashboards configured in your Uptime.com account. Dashboards provide customizable views of your monitoring data.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"dashboards": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all dashboards in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the dashboard",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the dashboard",
					},
					"ordering": schema.Int64Attribute{
						Computed:    true,
						Description: "Display order of the dashboard",
					},
					"is_pinned": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the dashboard is pinned",
					},
					"services_selected": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of selected service URLs",
					},
					"services_tags": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of service tags to filter by",
					},
					"metrics_show_section": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show metrics section",
					},
					"metrics_for_all_checks": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show metrics for all checks",
					},
					"services_show_section": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show services section",
					},
					"services_num_to_show": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of services to display",
					},
					"services_include_up": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to include UP services",
					},
					"services_include_down": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to include DOWN services",
					},
					"services_include_paused": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to include PAUSED services",
					},
					"services_include_maintenance": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to include services in MAINTENANCE",
					},
					"services_primary_sort": schema.StringAttribute{
						Computed:    true,
						Description: "Primary sort field for services",
					},
					"services_secondary_sort": schema.StringAttribute{
						Computed:    true,
						Description: "Secondary sort field for services",
					},
					"services_show_uptime": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show uptime for services",
					},
					"services_show_response_time": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show response time for services",
					},
					"alerts_show_section": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show alerts section",
					},
					"alerts_for_all_checks": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show alerts for all checks",
					},
				},
			},
		},
	},
}

type DashboardsDataSourceModel struct {
	ID         types.String                    `tfsdk:"id"`
	Dashboards []DashboardsDataSourceItemModel `tfsdk:"dashboards"`
}

type DashboardsDataSourceItemModel struct {
	ID                         types.Int64  `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Ordering                   types.Int64  `tfsdk:"ordering"`
	IsPinned                   types.Bool   `tfsdk:"is_pinned"`
	ServicesSelected           types.List   `tfsdk:"services_selected"`
	ServicesTags               types.List   `tfsdk:"services_tags"`
	MetricsShowSection         types.Bool   `tfsdk:"metrics_show_section"`
	MetricsForAllChecks        types.Bool   `tfsdk:"metrics_for_all_checks"`
	ServicesShowSection        types.Bool   `tfsdk:"services_show_section"`
	ServicesNumToShow          types.Int64  `tfsdk:"services_num_to_show"`
	ServicesIncludeUp          types.Bool   `tfsdk:"services_include_up"`
	ServicesIncludeDown        types.Bool   `tfsdk:"services_include_down"`
	ServicesIncludePaused      types.Bool   `tfsdk:"services_include_paused"`
	ServicesIncludeMaintenance types.Bool   `tfsdk:"services_include_maintenance"`
	ServicesPrimarySort        types.String `tfsdk:"services_primary_sort"`
	ServicesSecondarySort      types.String `tfsdk:"services_secondary_sort"`
	ServicesShowUptime         types.Bool   `tfsdk:"services_show_uptime"`
	ServicesShowResponseTime   types.Bool   `tfsdk:"services_show_response_time"`
	AlertsShowSection          types.Bool   `tfsdk:"alerts_show_section"`
	AlertsForAllChecks         types.Bool   `tfsdk:"alerts_for_all_checks"`
}

var _ datasource.DataSource = &DashboardsDataSource{}

type DashboardsDataSource struct {
	p *providerImpl
}

func (d DashboardsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_dashboards"
}

func (d DashboardsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = DashboardsDataSchema
}

func (d DashboardsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.Dashboards().List(ctx, upapi.DashboardListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := DashboardsDataSourceModel{
		ID:         types.StringValue(""),
		Dashboards: make([]DashboardsDataSourceItemModel, len(api.Items)),
	}

	for i := range api.Items {
		// Convert ServicesSelected slice to types.List
		servicesSelectedValues := make([]attr.Value, len(api.Items[i].ServicesSelected))
		for j, v := range api.Items[i].ServicesSelected {
			servicesSelectedValues[j] = types.StringValue(v)
		}
		servicesSelected := types.ListNull(types.StringType)
		if len(servicesSelectedValues) > 0 {
			servicesSelected = types.ListValueMust(types.StringType, servicesSelectedValues)
		}

		// Convert ServicesTags slice to types.List
		servicesTagsValues := make([]attr.Value, len(api.Items[i].ServicesTags))
		for j, v := range api.Items[i].ServicesTags {
			servicesTagsValues[j] = types.StringValue(v)
		}
		servicesTags := types.ListNull(types.StringType)
		if len(servicesTagsValues) > 0 {
			servicesTags = types.ListValueMust(types.StringType, servicesTagsValues)
		}

		model.Dashboards[i] = DashboardsDataSourceItemModel{
			ID:                         types.Int64Value(api.Items[i].PK),
			Name:                       types.StringValue(api.Items[i].Name),
			Ordering:                   types.Int64Value(api.Items[i].Ordering),
			IsPinned:                   types.BoolValue(api.Items[i].IsPinned),
			ServicesSelected:           servicesSelected,
			ServicesTags:               servicesTags,
			MetricsShowSection:         types.BoolValue(api.Items[i].MetricsShowSection),
			MetricsForAllChecks:        types.BoolValue(api.Items[i].MetricsForAllChecks),
			ServicesShowSection:        types.BoolValue(api.Items[i].ServicesShowSection),
			ServicesNumToShow:          types.Int64Value(api.Items[i].ServicesNumToShow),
			ServicesIncludeUp:          types.BoolValue(api.Items[i].ServicesIncludeUp),
			ServicesIncludeDown:        types.BoolValue(api.Items[i].ServicesIncludeDown),
			ServicesIncludePaused:      types.BoolValue(api.Items[i].ServicesIncludePaused),
			ServicesIncludeMaintenance: types.BoolValue(api.Items[i].ServicesIncludeMaintenance),
			ServicesPrimarySort:        types.StringValue(api.Items[i].ServicesPrimarySort),
			ServicesSecondarySort:      types.StringValue(api.Items[i].ServicesSecondarySort),
			ServicesShowUptime:         types.BoolValue(api.Items[i].ServicesShowUptime),
			ServicesShowResponseTime:   types.BoolValue(api.Items[i].ServicesShowResponseTime),
			AlertsShowSection:          types.BoolValue(api.Items[i].AlertsShowSection),
			AlertsForAllChecks:         types.BoolValue(api.Items[i].AlertsForAllChecks),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
