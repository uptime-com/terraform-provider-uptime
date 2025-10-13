package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageComponentDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageComponentDataSource{p: p}
}

var StatusPageComponentDataSchema = schema.Schema{
	Description: "Retrieve a list of all components for a specific status page.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspage_id": schema.Int64Attribute{
			Required:    true,
			Description: "ID of the status page to retrieve components for",
		},
		"components": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all components for the status page",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the component",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the component",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the component",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Description of the component",
					},
					"is_group": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this component is a group",
					},
					"group_id": schema.Int64Attribute{
						Computed:    true,
						Description: "ID of the parent group, if any",
					},
					"service_id": schema.Int64Attribute{
						Computed:    true,
						Description: "ID of the associated service/check",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "Current status of the component",
					},
					"auto_status_down": schema.StringAttribute{
						Computed:    true,
						Description: "Status to set when service is down",
					},
					"auto_status_up": schema.StringAttribute{
						Computed:    true,
						Description: "Status to set when service is up",
					},
				},
			},
		},
	},
}

type StatusPageComponentDataSourceModel struct {
	ID           types.String                             `tfsdk:"id"`
	StatusPageID types.Int64                              `tfsdk:"statuspage_id"`
	Components   []StatusPageComponentDataSourceItemModel `tfsdk:"components"`
}

type StatusPageComponentDataSourceItemModel struct {
	ID             types.Int64  `tfsdk:"id"`
	URL            types.String `tfsdk:"url"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	IsGroup        types.Bool   `tfsdk:"is_group"`
	GroupID        types.Int64  `tfsdk:"group_id"`
	ServiceID      types.Int64  `tfsdk:"service_id"`
	Status         types.String `tfsdk:"status"`
	AutoStatusDown types.String `tfsdk:"auto_status_down"`
	AutoStatusUp   types.String `tfsdk:"auto_status_up"`
}

var _ datasource.DataSource = &StatusPageComponentDataSource{}

type StatusPageComponentDataSource struct {
	p *providerImpl
}

func (d StatusPageComponentDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspage_components"
}

func (d StatusPageComponentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageComponentDataSchema
}

func (d StatusPageComponentDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config StatusPageComponentDataSourceModel
	diags := rq.Config.Get(ctx, &config)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	pk := upapi.PrimaryKey(config.StatusPageID.ValueInt64())
	api, err := d.p.api.StatusPages().Components(pk).List(ctx, upapi.StatusPageComponentListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageComponentDataSourceModel{
		ID:           types.StringValue(""),
		StatusPageID: config.StatusPageID,
		Components:   make([]StatusPageComponentDataSourceItemModel, len(api)),
	}

	for i := range api {
		groupID := types.Int64Null()
		if api[i].GroupID != nil {
			groupID = types.Int64Value(*api[i].GroupID)
		}

		serviceID := types.Int64Null()
		if api[i].ServiceID != nil {
			serviceID = types.Int64Value(*api[i].ServiceID)
		}

		model.Components[i] = StatusPageComponentDataSourceItemModel{
			ID:             types.Int64Value(api[i].PK),
			URL:            types.StringValue(api[i].URL),
			Name:           types.StringValue(api[i].Name),
			Description:    types.StringValue(api[i].Description),
			IsGroup:        types.BoolValue(api[i].IsGroup),
			GroupID:        groupID,
			ServiceID:      serviceID,
			Status:         types.StringValue(api[i].Status),
			AutoStatusDown: types.StringValue(api[i].AutoStatusDown),
			AutoStatusUp:   types.StringValue(api[i].AutoStatusUp),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
