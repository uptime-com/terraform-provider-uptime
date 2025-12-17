package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewSLAReportsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return SLAReportsDataSource{p: p}
}

// SLAReportsDataSchema defines the schema for the SLA reports data source.
var SLAReportsDataSchema = schema.Schema{
	Description: "Retrieve a list of all SLA reports configured in your Uptime.com account. SLA reports track service level agreement compliance.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"sla_reports": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all SLA reports in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the SLA report",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the SLA report resource",
					},
					"stats_url": schema.StringAttribute{
						Computed:    true,
						Description: "URL for the SLA report statistics",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the SLA report",
					},
					"services_tags": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of service tags to filter by",
					},
					"default_date_range": schema.StringAttribute{
						Computed:    true,
						Description: "Default date range for the report",
					},
					"filter_with_downtime": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to filter services with downtime",
					},
					"filter_uptime_sla_violations": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to filter services with uptime SLA violations",
					},
					"filter_slowest": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to filter slowest services",
					},
					"filter_response_time_sla_violations": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to filter services with response time SLA violations",
					},
					"show_uptime_section": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show uptime section in the report",
					},
					"show_uptime_sla": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show uptime SLA in the report",
					},
					"show_response_time_section": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to show response time section in the report",
					},
				},
			},
		},
	},
}

type SLAReportsDataSourceModel struct {
	ID         types.String                    `tfsdk:"id"`
	SLAReports []SLAReportsDataSourceItemModel `tfsdk:"sla_reports"`
}

type SLAReportsDataSourceItemModel struct {
	ID                              types.Int64  `tfsdk:"id"`
	URL                             types.String `tfsdk:"url"`
	StatsURL                        types.String `tfsdk:"stats_url"`
	Name                            types.String `tfsdk:"name"`
	ServicesTags                    types.List   `tfsdk:"services_tags"`
	DefaultDateRange                types.String `tfsdk:"default_date_range"`
	FilterWithDowntime              types.Bool   `tfsdk:"filter_with_downtime"`
	FilterUptimeSLAViolations       types.Bool   `tfsdk:"filter_uptime_sla_violations"`
	FilterSlowest                   types.Bool   `tfsdk:"filter_slowest"`
	FilterResponseTimeSLAViolations types.Bool   `tfsdk:"filter_response_time_sla_violations"`
	ShowUptimeSection               types.Bool   `tfsdk:"show_uptime_section"`
	ShowUptimeSLA                   types.Bool   `tfsdk:"show_uptime_sla"`
	ShowResponseTimeSection         types.Bool   `tfsdk:"show_response_time_section"`
}

var _ datasource.DataSource = &SLAReportsDataSource{}

type SLAReportsDataSource struct {
	p *providerImpl
}

func (d SLAReportsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_sla_reports"
}

func (d SLAReportsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = SLAReportsDataSchema
}

func (d SLAReportsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.SLAReports().List(ctx, upapi.SLAReportListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := SLAReportsDataSourceModel{
		ID:         types.StringValue(""),
		SLAReports: make([]SLAReportsDataSourceItemModel, len(api.Items)),
	}

	for i := range api.Items {
		// Convert ServicesTags slice to types.List
		servicesTagsValues := make([]attr.Value, len(api.Items[i].ServicesTags))
		for j, v := range api.Items[i].ServicesTags {
			servicesTagsValues[j] = types.StringValue(v)
		}
		servicesTags := types.ListNull(types.StringType)
		if len(servicesTagsValues) > 0 {
			servicesTags = types.ListValueMust(types.StringType, servicesTagsValues)
		}

		model.SLAReports[i] = SLAReportsDataSourceItemModel{
			ID:                              types.Int64Value(api.Items[i].PK),
			URL:                             types.StringValue(api.Items[i].URL),
			StatsURL:                        types.StringValue(api.Items[i].StatsURL),
			Name:                            types.StringValue(api.Items[i].Name),
			ServicesTags:                    servicesTags,
			DefaultDateRange:                types.StringValue(api.Items[i].DefaultDateRange),
			FilterWithDowntime:              types.BoolValue(api.Items[i].FilterWithDowntime),
			FilterUptimeSLAViolations:       types.BoolValue(api.Items[i].FilterUptimeSLAViolations),
			FilterSlowest:                   types.BoolValue(api.Items[i].FilterSlowest),
			FilterResponseTimeSLAViolations: types.BoolValue(api.Items[i].FilterResponseTimeSLAViolations),
			ShowUptimeSection:               types.BoolValue(api.Items[i].ShowUptimeSection),
			ShowUptimeSLA:                   types.BoolValue(api.Items[i].ShowUptimeSLA),
			ShowResponseTimeSection:         types.BoolValue(api.Items[i].ShowResponseTimeSection),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
