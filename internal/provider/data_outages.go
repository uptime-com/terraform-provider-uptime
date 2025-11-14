package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewOutagesDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return OutagesDataSource{p: p}
}

// OutagesDataSchema defines the schema for the outages data source.
var OutagesDataSchema = schema.Schema{
	Description: "Retrieve a list of outages from your Uptime.com account. Outages represent periods when a check was down.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"outages": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of outages",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the outage",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the outage resource",
					},
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the outage started",
					},
					"resolved_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the outage was resolved",
					},
					"duration_secs": schema.Int64Attribute{
						Computed:    true,
						Description: "Duration of the outage in seconds",
					},
					"ignore_alert_url": schema.StringAttribute{
						Computed:    true,
						Description: "URL to ignore this outage alert",
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
					"state_is_up": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the check state is currently UP",
					},
					"ignored": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the outage is ignored",
					},
					"num_locations_down": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of monitoring locations that detected the outage",
					},
				},
			},
		},
	},
}

type OutagesDataSourceModel struct {
	ID      types.String                 `tfsdk:"id"`
	Outages []OutagesDataSourceItemModel `tfsdk:"outages"`
}

type OutagesDataSourceItemModel struct {
	ID                         types.Int64  `tfsdk:"id"`
	URL                        types.String `tfsdk:"url"`
	CreatedAt                  types.String `tfsdk:"created_at"`
	ResolvedAt                 types.String `tfsdk:"resolved_at"`
	DurationSecs               types.Int64  `tfsdk:"duration_secs"`
	IgnoreAlertURL             types.String `tfsdk:"ignore_alert_url"`
	CheckPK                    types.Int64  `tfsdk:"check_pk"`
	CheckURL                   types.String `tfsdk:"check_url"`
	CheckAddress               types.String `tfsdk:"check_address"`
	CheckName                  types.String `tfsdk:"check_name"`
	CheckMonitoringServiceType types.String `tfsdk:"check_monitoring_service_type"`
	StateIsUp                  types.Bool   `tfsdk:"state_is_up"`
	Ignored                    types.Bool   `tfsdk:"ignored"`
	NumLocationsDown           types.Int64  `tfsdk:"num_locations_down"`
}

var _ datasource.DataSource = &OutagesDataSource{}

type OutagesDataSource struct {
	p *providerImpl
}

func (d OutagesDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_outages"
}

func (d OutagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = OutagesDataSchema
}

func (d OutagesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.Outages().List(ctx, upapi.OutageListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := OutagesDataSourceModel{
		ID:      types.StringValue(""),
		Outages: make([]OutagesDataSourceItemModel, len(api)),
	}

	for i := range api {
		createdAt := ""
		if !api[i].CreatedAt.IsZero() {
			createdAt = api[i].CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		resolvedAt := ""
		if !api[i].ResolvedAt.IsZero() {
			resolvedAt = api[i].ResolvedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		model.Outages[i] = OutagesDataSourceItemModel{
			ID:                         types.Int64Value(api[i].PK),
			URL:                        types.StringValue(api[i].URL),
			CreatedAt:                  types.StringValue(createdAt),
			ResolvedAt:                 types.StringValue(resolvedAt),
			DurationSecs:               types.Int64Value(api[i].DurationSecs),
			IgnoreAlertURL:             types.StringValue(api[i].IgnoreAlertURL),
			CheckPK:                    types.Int64Value(api[i].CheckPK),
			CheckURL:                   types.StringValue(api[i].CheckURL),
			CheckAddress:               types.StringValue(api[i].CheckAddress),
			CheckName:                  types.StringValue(api[i].CheckName),
			CheckMonitoringServiceType: types.StringValue(api[i].CheckMonitoringServiceType),
			StateIsUp:                  types.BoolValue(api[i].StateIsUp),
			Ignored:                    types.BoolValue(api[i].Ignored),
			NumLocationsDown:           types.Int64Value(api[i].NumLocationsDown),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
