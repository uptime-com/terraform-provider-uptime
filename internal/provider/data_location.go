package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewLocationsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return &locationsDataSource{p: p}
}

var locationsDataSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"locations": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"location": schema.StringAttribute{
						Computed: true,
					},
					"ip": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	},
}

type locationsDataSourceLocationModel struct {
	ID       types.String `tfsdk:"id"       ref:",skip"`
	Name     types.String `tfsdk:"name"     ref:"ProbeName"`
	Location types.String `tfsdk:"location"`
	IP       types.String `tfsdk:"ip"       ref:"IPAddress"`
}

type locationsDataSourceModel struct {
	ID        types.String                       `tfsdk:"id"`
	Locations []locationsDataSourceLocationModel `tfsdk:"locations"`
}

var _ datasource.DataSource = &locationsDataSource{}

type locationsDataSource struct {
	p *providerImpl
}

func (d *locationsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_locations"
}

func (d *locationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = locationsDataSchema
}

func (d *locationsDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	apires, err := d.p.api.ProbeServers().List(ctx)
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}
	state := locationsDataSourceModel{
		ID:        types.StringValue(""),
		Locations: make([]locationsDataSourceLocationModel, len(apires)),
	}
	for i := range apires {
		diags := valueFromAPI(&state.Locations[i], apires[i])
		if diags.HasError() {
			rs.Diagnostics = diags
			return
		}
	}
	diags := rs.State.Set(ctx, &state)
	if diags.HasError() {
		rs.Diagnostics = diags
		return
	}
	return
}
