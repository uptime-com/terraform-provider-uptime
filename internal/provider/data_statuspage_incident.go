package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageIncidentDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageIncidentDataSource{p: p}
}

var StatusPageIncidentDataSchema = schema.Schema{
	Description: "Retrieve a list of all incidents for a specific status page.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspage_id": schema.Int64Attribute{
			Required:    true,
			Description: "ID of the status page to retrieve incidents for",
		},
		"incidents": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all incidents for the status page",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the incident",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the incident",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the incident",
					},
					"incident_type": schema.StringAttribute{
						Computed:    true,
						Description: "Type of incident",
					},
					"starts_at": schema.StringAttribute{
						Computed:    true,
						Description: "Start time of the incident",
					},
					"ends_at": schema.StringAttribute{
						Computed:    true,
						Description: "End time of the incident",
					},
					"include_in_global_metrics": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to include in global metrics",
					},
				},
			},
		},
	},
}

type StatusPageIncidentDataSourceModel struct {
	ID           types.String                            `tfsdk:"id"`
	StatusPageID types.Int64                             `tfsdk:"statuspage_id"`
	Incidents    []StatusPageIncidentDataSourceItemModel `tfsdk:"incidents"`
}

type StatusPageIncidentDataSourceItemModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	URL                    types.String `tfsdk:"url"`
	Name                   types.String `tfsdk:"name"`
	IncidentType           types.String `tfsdk:"incident_type"`
	StartsAt               types.String `tfsdk:"starts_at"`
	EndsAt                 types.String `tfsdk:"ends_at"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
}

var _ datasource.DataSource = &StatusPageIncidentDataSource{}

type StatusPageIncidentDataSource struct {
	p *providerImpl
}

func (d StatusPageIncidentDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspage_incidents"
}

func (d StatusPageIncidentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageIncidentDataSchema
}

func (d StatusPageIncidentDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config StatusPageIncidentDataSourceModel
	diags := rq.Config.Get(ctx, &config)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	pk := upapi.PrimaryKey(config.StatusPageID.ValueInt64())
	api, err := d.p.api.StatusPages().Incidents(pk).List(ctx, upapi.StatusPageIncidentListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageIncidentDataSourceModel{
		ID:           types.StringValue(""),
		StatusPageID: config.StatusPageID,
		Incidents:    make([]StatusPageIncidentDataSourceItemModel, len(api)),
	}

	for i := range api {
		model.Incidents[i] = StatusPageIncidentDataSourceItemModel{
			ID:                     types.Int64Value(api[i].PK),
			URL:                    types.StringValue(api[i].URL),
			Name:                   types.StringValue(api[i].Name),
			IncidentType:           types.StringValue(api[i].IncidentType),
			StartsAt:               types.StringValue(api[i].StartsAt),
			EndsAt:                 types.StringValue(api[i].EndsAt),
			IncludeInGlobalMetrics: types.BoolValue(api[i].IncludeInGlobalMetrics),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
