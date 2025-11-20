package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageCurrentStatusDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageCurrentStatusDataSource{p: p}
}

var StatusPageCurrentStatusDataSchema = schema.Schema{
	Description: "Retrieve the current operational status of a status page including component statuses.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspage_id": schema.Int64Attribute{
			Required:    true,
			Description: "ID of the status page to retrieve current status for",
		},
		"global_is_operational": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether the status page is operational (all components operational)",
		},
		"components": schema.ListNestedAttribute{
			Computed:    true,
			Description: "Current status of all components",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Component ID",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Component name",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "Current status of the component",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Status description",
					},
				},
			},
		},
	},
}

type StatusPageCurrentStatusDataSourceModel struct {
	ID                  types.String                                      `tfsdk:"id"`
	StatusPageID        types.Int64                                       `tfsdk:"statuspage_id"`
	GlobalIsOperational types.Bool                                        `tfsdk:"global_is_operational"`
	Components          []StatusPageCurrentStatusDataSourceComponentModel `tfsdk:"components"`
}

type StatusPageCurrentStatusDataSourceComponentModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
	Description types.String `tfsdk:"description"`
}

var _ datasource.DataSource = &StatusPageCurrentStatusDataSource{}

type StatusPageCurrentStatusDataSource struct {
	p *providerImpl
}

func (d StatusPageCurrentStatusDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspage_current_status"
}

func (d StatusPageCurrentStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageCurrentStatusDataSchema
}

func (d StatusPageCurrentStatusDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config StatusPageCurrentStatusDataSourceModel
	diags := rq.Config.Get(ctx, &config)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	pk := upapi.PrimaryKey(config.StatusPageID.ValueInt64())
	api, err := d.p.api.StatusPages().CurrentStatus(pk).Get(ctx)
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageCurrentStatusDataSourceModel{
		ID:                  types.StringValue(""),
		StatusPageID:        config.StatusPageID,
		GlobalIsOperational: types.BoolValue(api.GlobalIsOperational),
		Components:          make([]StatusPageCurrentStatusDataSourceComponentModel, len(api.Components)),
	}

	for i := range api.Components {
		model.Components[i] = StatusPageCurrentStatusDataSourceComponentModel{
			ID:          types.Int64Value(api.Components[i].PK),
			Name:        types.StringValue(api.Components[i].Name),
			Status:      types.StringValue(api.Components[i].Status),
			Description: types.StringValue(api.Components[i].Description),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
