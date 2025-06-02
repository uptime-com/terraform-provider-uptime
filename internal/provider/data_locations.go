package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
					"ipv6": schema.StringAttribute{
						Computed: true,
					},
					"ipv4_addresses": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"ipv6_addresses": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
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
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Location      types.String `tfsdk:"location"`
	IP            types.String `tfsdk:"ip"`
	IPv6          types.String `tfsdk:"ipv6"`
	IPv4Addresses types.List   `tfsdk:"ipv4_addresses"`
	IPv6Addresses types.List   `tfsdk:"ipv6_addresses"`
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
		var primaryIP string
		var primaryIPv6 string

		// Get primary IPs (first in each list)
		if len(api[i].IPv4Addresses) > 0 {
			primaryIP = api[i].IPv4Addresses[0]
		}
		if len(api[i].IPv6Addresses) > 0 {
			primaryIPv6 = api[i].IPv6Addresses[0]
		}

		// Convert to Terraform types
		ipv4StringList := make([]types.String, len(api[i].IPv4Addresses))
		for j, ip := range api[i].IPv4Addresses {
			ipv4StringList[j] = types.StringValue(ip)
		}

		ipv6StringList := make([]types.String, len(api[i].IPv6Addresses))
		for j, ip := range api[i].IPv6Addresses {
			ipv6StringList[j] = types.StringValue(ip)
		}

		model.Locations[i] = LocationsDataSourceLocationModel{
			Name:          types.StringValue(api[i].ProbeName),
			Location:      types.StringValue(api[i].Location),
			IP:            types.StringValue(primaryIP),
			IPv6:          types.StringValue(primaryIPv6),
			IPv4Addresses: types.ListValueMust(types.StringType, convertToAttrValues(ipv4StringList)),
			IPv6Addresses: types.ListValueMust(types.StringType, convertToAttrValues(ipv6StringList)),
		}
	}
	rs.Diagnostics = rs.State.Set(ctx, model)
	return
}

func convertToAttrValues(stringValues []types.String) []attr.Value {
	result := make([]attr.Value, len(stringValues))
	for i, v := range stringValues {
		result[i] = v
	}
	return result
}
