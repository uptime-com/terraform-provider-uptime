package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageStatusHistoryDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageStatusHistoryDataSource{p: p}
}

var StatusPageStatusHistoryDataSchema = schema.Schema{
	Description: "Retrieve status history entries for a status page.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspage_id": schema.Int64Attribute{
			Required:    true,
			Description: "ID of the status page to retrieve status history for",
		},
		"status": schema.StringAttribute{
			Optional:    true,
			Description: "Filter by status",
		},
		"component_id": schema.Int64Attribute{
			Optional:    true,
			Description: "Filter by component ID",
		},
		"date_from": schema.StringAttribute{
			Optional:    true,
			Description: "Filter entries from this date (ISO 8601 format)",
		},
		"date_to": schema.StringAttribute{
			Optional:    true,
			Description: "Filter entries until this date (ISO 8601 format)",
		},
		"history": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of status history entries",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "History entry ID",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "Status value",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Status description",
					},
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "When this status was created",
					},
					"updated_at": schema.StringAttribute{
						Computed:    true,
						Description: "When this status was last updated",
					},
					"component_id": schema.Int64Attribute{
						Computed:    true,
						Description: "Associated component ID (if applicable)",
					},
				},
			},
		},
	},
}

type StatusPageStatusHistoryDataSourceModel struct {
	ID           types.String                                 `tfsdk:"id"`
	StatusPageID types.Int64                                  `tfsdk:"statuspage_id"`
	Status       types.String                                 `tfsdk:"status"`
	ComponentID  types.Int64                                  `tfsdk:"component_id"`
	DateFrom     types.String                                 `tfsdk:"date_from"`
	DateTo       types.String                                 `tfsdk:"date_to"`
	History      []StatusPageStatusHistoryDataSourceItemModel `tfsdk:"history"`
}

type StatusPageStatusHistoryDataSourceItemModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Status      types.String `tfsdk:"status"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	ComponentID types.Int64  `tfsdk:"component_id"`
}

var _ datasource.DataSource = &StatusPageStatusHistoryDataSource{}

type StatusPageStatusHistoryDataSource struct {
	p *providerImpl
}

func (d StatusPageStatusHistoryDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspage_status_history"
}

func (d StatusPageStatusHistoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageStatusHistoryDataSchema
}

func (d StatusPageStatusHistoryDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config StatusPageStatusHistoryDataSourceModel
	diags := rq.Config.Get(ctx, &config)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	pk := upapi.PrimaryKey(config.StatusPageID.ValueInt64())
	opts := upapi.StatusPageStatusHistoryListOptions{}

	if !config.Status.IsNull() {
		opts.Status = config.Status.ValueString()
	}
	if !config.ComponentID.IsNull() {
		opts.ComponentPK = config.ComponentID.ValueInt64()
	}
	if !config.DateFrom.IsNull() {
		opts.DateFrom = config.DateFrom.ValueString()
	}
	if !config.DateTo.IsNull() {
		opts.DateTo = config.DateTo.ValueString()
	}

	api, err := d.p.api.StatusPages().StatusHistory(pk).List(ctx, opts)
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageStatusHistoryDataSourceModel{
		ID:           types.StringValue(""),
		StatusPageID: config.StatusPageID,
		Status:       config.Status,
		ComponentID:  config.ComponentID,
		DateFrom:     config.DateFrom,
		DateTo:       config.DateTo,
		History:      make([]StatusPageStatusHistoryDataSourceItemModel, len(api.Items)),
	}

	for i := range api.Items {
		componentID := types.Int64Null()
		if api.Items[i].ComponentPK != nil {
			componentID = types.Int64Value(*api.Items[i].ComponentPK)
		}

		model.History[i] = StatusPageStatusHistoryDataSourceItemModel{
			ID:          types.Int64Value(api.Items[i].PK),
			Status:      types.StringValue(api.Items[i].Status),
			Description: types.StringValue(api.Items[i].Description),
			CreatedAt:   types.StringValue(api.Items[i].CreatedAt),
			UpdatedAt:   types.StringValue(api.Items[i].UpdatedAt),
			ComponentID: componentID,
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
