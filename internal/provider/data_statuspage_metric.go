package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageMetricDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageMetricDataSource{p: p}
}

var StatusPageMetricDataSchema = schema.Schema{
	Description: "Retrieve a list of all metrics for a specific status page.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspage_id": schema.Int64Attribute{
			Required:    true,
			Description: "ID of the status page to retrieve metrics for",
		},
		"metrics": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all metrics for the status page",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the metric",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the metric",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the metric",
					},
					"service_id": schema.Int64Attribute{
						Computed:    true,
						Description: "ID of the associated service/check",
					},
					"is_visible": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the metric is visible on the status page",
					},
				},
			},
		},
	},
}

type StatusPageMetricDataSourceModel struct {
	ID           types.String                          `tfsdk:"id"`
	StatusPageID types.Int64                           `tfsdk:"statuspage_id"`
	Metrics      []StatusPageMetricDataSourceItemModel `tfsdk:"metrics"`
}

type StatusPageMetricDataSourceItemModel struct {
	ID        types.Int64  `tfsdk:"id"`
	URL       types.String `tfsdk:"url"`
	Name      types.String `tfsdk:"name"`
	ServiceID types.Int64  `tfsdk:"service_id"`
	IsVisible types.Bool   `tfsdk:"is_visible"`
}

var _ datasource.DataSource = &StatusPageMetricDataSource{}

type StatusPageMetricDataSource struct {
	p *providerImpl
}

func (d StatusPageMetricDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspage_metrics"
}

func (d StatusPageMetricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageMetricDataSchema
}

func (d StatusPageMetricDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config StatusPageMetricDataSourceModel
	diags := rq.Config.Get(ctx, &config)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	pk := upapi.PrimaryKey(config.StatusPageID.ValueInt64())
	api, err := d.p.api.StatusPages().Metrics(pk).List(ctx, upapi.StatusPageMetricListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageMetricDataSourceModel{
		ID:           types.StringValue(""),
		StatusPageID: config.StatusPageID,
		Metrics:      make([]StatusPageMetricDataSourceItemModel, len(api)),
	}

	for i := range api {
		model.Metrics[i] = StatusPageMetricDataSourceItemModel{
			ID:        types.Int64Value(api[i].PK),
			URL:       types.StringValue(api[i].URL),
			Name:      types.StringValue(api[i].Name),
			ServiceID: types.Int64Value(api[i].ServiceID),
			IsVisible: types.BoolValue(api[i].IsVisible),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
