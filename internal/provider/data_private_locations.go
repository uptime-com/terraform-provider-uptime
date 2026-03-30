package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// extendedProbeServer extends the upstream client's ProbeServer with fields
// that the API returns but the client struct does not yet include.
type extendedProbeServer struct {
	Location      string   `json:"location"`
	Country       string   `json:"country"`
	ProbeName     string   `json:"probe_name"`
	IPAddress     string   `json:"ip_address"`
	IPv4Addresses []string `json:"ipv4_addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
	IsPrivate     bool     `json:"is_private"`
}

// listExtendedProbeServers fetches probe servers directly from the API,
// deserializing into extendedProbeServer to capture fields like is_private
// that the upstream Go client does not yet include.
func (p *providerImpl) listExtendedProbeServers(ctx context.Context) ([]extendedProbeServer, error) {
	url := strings.TrimRight(p.baseURL, "/") + "/probe-servers/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+p.token)
	req.Header.Set("User-Agent", p.UserAgentString())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var servers []extendedProbeServer
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return servers, nil
}

func NewPrivateLocationsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return PrivateLocationsDataSource{p: p}
}

var PrivateLocationsDataSchema = schema.Schema{
	Description: "List private monitoring locations available in the account.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"locations": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Probe server name.",
					},
					"location": schema.StringAttribute{
						Computed:    true,
						Description: "Location identifier, usable in check location attributes.",
					},
					"country": schema.StringAttribute{
						Computed:    true,
						Description: "Country where the private location is deployed.",
					},
					"ip": schema.StringAttribute{
						Computed:    true,
						Description: "Primary IPv4 address.",
					},
					"ipv6": schema.StringAttribute{
						Computed:    true,
						Description: "Primary IPv6 address.",
					},
					"ipv4_addresses": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "All IPv4 addresses.",
					},
					"ipv6_addresses": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "All IPv6 addresses.",
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
	servers, err := d.p.listExtendedProbeServers(ctx)
	if err != nil {
		rs.Diagnostics.AddError("Failed to list probe servers", err.Error())
		return
	}

	var privateServers []extendedProbeServer
	for _, s := range servers {
		if s.IsPrivate {
			privateServers = append(privateServers, s)
		}
	}

	model := PrivateLocationsDataSourceModel{
		ID:        types.StringValue(""),
		Locations: make([]PrivateLocationsDataSourceLocationModel, len(privateServers)),
	}
	for i, s := range privateServers {
		var primaryIP, primaryIPv6 string
		if len(s.IPv4Addresses) > 0 {
			primaryIP = s.IPv4Addresses[0]
		}
		if len(s.IPv6Addresses) > 0 {
			primaryIPv6 = s.IPv6Addresses[0]
		}

		ipv4Values := make([]attr.Value, len(s.IPv4Addresses))
		for j, ip := range s.IPv4Addresses {
			ipv4Values[j] = types.StringValue(ip)
		}
		ipv6Values := make([]attr.Value, len(s.IPv6Addresses))
		for j, ip := range s.IPv6Addresses {
			ipv6Values[j] = types.StringValue(ip)
		}

		model.Locations[i] = PrivateLocationsDataSourceLocationModel{
			Name:          types.StringValue(s.ProbeName),
			Location:      types.StringValue(s.Location),
			Country:       types.StringValue(s.Country),
			IP:            types.StringValue(primaryIP),
			IPv6:          types.StringValue(primaryIPv6),
			IPv4Addresses: types.ListValueMust(types.StringType, ipv4Values),
			IPv6Addresses: types.ListValueMust(types.StringType, ipv6Values),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
