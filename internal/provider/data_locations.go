package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewLocationsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return LocationsDataSource{p: p}
}

var LocationsDataSchema = schema.Schema{
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

type LocationsDataSourceModel struct {
	ID        types.String                       `tfsdk:"id"`
	Locations []LocationsDataSourceLocationModel `tfsdk:"locations"`
}

type LocationsDataSourceLocationModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Location types.String `tfsdk:"location"`
	IP       types.String `tfsdk:"ip"`
}

var _ datasource.DataSource = &LocationsDataSource{}

type LocationsDataSource struct {
	p *providerImpl
}

func (d LocationsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_locations"
}

func (d LocationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = LocationsDataSchema
}

func (d LocationsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.ProbeServers().List(ctx)
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}
	model := LocationsDataSourceModel{
		ID:        types.StringValue(""),
		Locations: make([]LocationsDataSourceLocationModel, len(api)),
	}
	for i := range api {
		model.Locations[i] = LocationsDataSourceLocationModel{
			Name:     types.StringValue(api[i].ProbeName),
			Location: types.StringValue(api[i].Location),
			IP:       types.StringValue(api[i].IPAddress),
		}
	}
	rs.Diagnostics = rs.State.Set(ctx, model)
	return
}
