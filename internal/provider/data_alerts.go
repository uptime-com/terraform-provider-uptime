package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewAlertsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return AlertsDataSource{p: p}
}

// AlertsDataSchema defines the schema for the alerts data source.
var AlertsDataSchema = schema.Schema{
	Description: "Retrieve a list of alerts from your Uptime.com account. Alerts are generated when a check detects an issue from a monitoring location.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"alerts": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of alerts",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the alert",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the alert resource",
					},
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the alert was created",
					},
					"resolved_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the alert was resolved",
					},
					"monitoring_server_name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the monitoring server that detected the alert",
					},
					"location": schema.StringAttribute{
						Computed:    true,
						Description: "Location of the monitoring server",
					},
					"output": schema.StringAttribute{
						Computed:    true,
						Description: "Alert output message",
					},
					"state_is_up": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the check state is UP",
					},
					"ignored": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the alert is ignored",
					},
					"check_pk": schema.Int64Attribute{
						Computed:    true,
						Description: "Primary key of the associated check",
					},
					"check_url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL of the associated check",
					},
					"check_address": schema.StringAttribute{
						Computed:    true,
						Description: "Address of the check being monitored",
					},
					"check_name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the check",
					},
					"check_monitoring_service_type": schema.StringAttribute{
						Computed:    true,
						Description: "Type of monitoring service (e.g., 'HTTP', 'TCP')",
					},
				},
			},
		},
	},
}

type AlertsDataSourceModel struct {
	ID     types.String                `tfsdk:"id"`
	Alerts []AlertsDataSourceItemModel `tfsdk:"alerts"`
}

type AlertsDataSourceItemModel struct {
	ID                         types.Int64  `tfsdk:"id"`
	URL                        types.String `tfsdk:"url"`
	CreatedAt                  types.String `tfsdk:"created_at"`
	ResolvedAt                 types.String `tfsdk:"resolved_at"`
	MonitoringServerName       types.String `tfsdk:"monitoring_server_name"`
	Location                   types.String `tfsdk:"location"`
	Output                     types.String `tfsdk:"output"`
	StateIsUp                  types.Bool   `tfsdk:"state_is_up"`
	Ignored                    types.Bool   `tfsdk:"ignored"`
	CheckPK                    types.Int64  `tfsdk:"check_pk"`
	CheckURL                   types.String `tfsdk:"check_url"`
	CheckAddress               types.String `tfsdk:"check_address"`
	CheckName                  types.String `tfsdk:"check_name"`
	CheckMonitoringServiceType types.String `tfsdk:"check_monitoring_service_type"`
}

var _ datasource.DataSource = &AlertsDataSource{}

type AlertsDataSource struct {
	p *providerImpl
}

func (d AlertsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_alerts"
}

func (d AlertsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = AlertsDataSchema
}

func (d AlertsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.Alerts().List(ctx, upapi.AlertListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := AlertsDataSourceModel{
		ID:     types.StringValue(""),
		Alerts: make([]AlertsDataSourceItemModel, len(api)),
	}

	for i := range api {
		createdAt := ""
		if api[i].CreatedAt != nil {
			createdAt = api[i].CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		resolvedAt := ""
		if api[i].ResolvedAt != nil {
			resolvedAt = api[i].ResolvedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		model.Alerts[i] = AlertsDataSourceItemModel{
			ID:                         types.Int64Value(api[i].PK),
			URL:                        types.StringValue(api[i].URL),
			CreatedAt:                  types.StringValue(createdAt),
			ResolvedAt:                 types.StringValue(resolvedAt),
			MonitoringServerName:       types.StringValue(api[i].MonitoringServerName),
			Location:                   types.StringValue(api[i].Location),
			Output:                     types.StringValue(api[i].Output),
			StateIsUp:                  types.BoolValue(api[i].StateIsUp),
			Ignored:                    types.BoolValue(api[i].Ignored),
			CheckPK:                    types.Int64Value(api[i].CheckPK),
			CheckURL:                   types.StringValue(api[i].CheckURL),
			CheckAddress:               types.StringValue(api[i].CheckAddress),
			CheckName:                  types.StringValue(api[i].CheckName),
			CheckMonitoringServiceType: types.StringValue(api[i].CheckMonitoringServiceType),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
