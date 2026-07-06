package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewPrivateLocationsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return PrivateLocationsDataSource{p: p}
}

var PrivateLocationsDataSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"locations": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed: true,
					},
					"location": schema.StringAttribute{
						Computed: true,
					},
					"country": schema.StringAttribute{
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

type PrivateLocationsDataSourceModel struct {
	ID        types.String                              `tfsdk:"id"`
	Locations []PrivateLocationsDataSourceLocationModel `tfsdk:"locations"`
}

type PrivateLocationsDataSourceLocationModel struct {
	Name          types.String `tfsdk:"name"`
	Location      types.String `tfsdk:"location"`
	Country       types.String `tfsdk:"country"`
	IP            types.String `tfsdk:"ip"`
	IPv6          types.String `tfsdk:"ipv6"`
	IPv4Addresses types.List   `tfsdk:"ipv4_addresses"`
	IPv6Addresses types.List   `tfsdk:"ipv6_addresses"`
}

var _ datasource.DataSource = &PrivateLocationsDataSource{}

type PrivateLocationsDataSource struct {
	p *providerImpl
}

func (d PrivateLocationsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_private_locations"
}

func (d PrivateLocationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = PrivateLocationsDataSchema
}

func (d PrivateLocationsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.ProbeServers().List(ctx)
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}
	model := PrivateLocationsDataSourceModel{
		ID:        types.StringValue(""),
		Locations: make([]PrivateLocationsDataSourceLocationModel, 0, len(api.Items)),
	}
	for i := range api.Items {
		if !api.Items[i].IsPrivate {
			continue
		}

		var primaryIP string
		var primaryIPv6 string

		if len(api.Items[i].IPv4Addresses) > 0 {
			primaryIP = api.Items[i].IPv4Addresses[0]
		}
		if len(api.Items[i].IPv6Addresses) > 0 {
			primaryIPv6 = api.Items[i].IPv6Addresses[0]
		}

		ipv4StringList := make([]types.String, len(api.Items[i].IPv4Addresses))
		for j, ip := range api.Items[i].IPv4Addresses {
			ipv4StringList[j] = types.StringValue(ip)
		}

		ipv6StringList := make([]types.String, len(api.Items[i].IPv6Addresses))
		for j, ip := range api.Items[i].IPv6Addresses {
			ipv6StringList[j] = types.StringValue(ip)
		}

		model.Locations = append(model.Locations, PrivateLocationsDataSourceLocationModel{
			Name:          types.StringValue(api.Items[i].ProbeName),
			Location:      types.StringValue(api.Items[i].Location),
			Country:       types.StringValue(api.Items[i].Country),
			IP:            types.StringValue(primaryIP),
			IPv6:          types.StringValue(primaryIPv6),
			IPv4Addresses: types.ListValueMust(types.StringType, convertToAttrValues(ipv4StringList)),
			IPv6Addresses: types.ListValueMust(types.StringType, convertToAttrValues(ipv6StringList)),
		})
	}
	rs.Diagnostics = rs.State.Set(ctx, model)
	return
}
