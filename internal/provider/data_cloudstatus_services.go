package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCloudStatusServicesDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return CloudStatusServicesDataSource{p: p}
}

var CloudStatusServicesDataSchema = schema.Schema{
	Description: "Look up Cloud Status services within a provider group. Use the returned `id` values to populate `services` on a `uptime_check_cloudstatus` resource.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"group": schema.StringAttribute{
			Optional:    true,
			Description: "Provider group ID or case-insensitive name substring. Strongly recommended because the services list is large; without it you will fetch every service across every provider.",
		},
		"search": schema.StringAttribute{
			Optional:    true,
			Description: "Case-insensitive substring matched against service name, title, or sub-title",
		},
		"services": schema.ListNestedAttribute{
			Computed:    true,
			Description: "Matching Cloud Status services",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the service",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Internal name (slug) for the service",
					},
					"title": schema.StringAttribute{
						Computed:    true,
						Description: "Human-readable service title",
					},
					"sub_title": schema.StringAttribute{
						Computed:    true,
						Description: "Optional secondary title (e.g., region/zone)",
					},
					"group_id": schema.Int64Attribute{
						Computed:    true,
						Description: "ID of the provider group this service belongs to",
					},
					"group": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the provider group this service belongs to",
					},
				},
			},
		},
	},
}

type CloudStatusServicesDataSourceModel struct {
	ID       types.String                             `tfsdk:"id"`
	Group    types.String                             `tfsdk:"group"`
	Search   types.String                             `tfsdk:"search"`
	Services []CloudStatusServicesDataSourceItemModel `tfsdk:"services"`
}

type CloudStatusServicesDataSourceItemModel struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Title    types.String `tfsdk:"title"`
	SubTitle types.String `tfsdk:"sub_title"`
	GroupID  types.Int64  `tfsdk:"group_id"`
	Group    types.String `tfsdk:"group"`
}

var _ datasource.DataSource = &CloudStatusServicesDataSource{}

type CloudStatusServicesDataSource struct {
	p *providerImpl
}

func (d CloudStatusServicesDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_cloudstatus_services"
}

func (d CloudStatusServicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = CloudStatusServicesDataSchema
}

func (d CloudStatusServicesDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config CloudStatusServicesDataSourceModel
	rs.Diagnostics.Append(rq.Config.Get(ctx, &config)...)
	if rs.Diagnostics.HasError() {
		return
	}

	const pageSize int64 = 100
	const maxPages int64 = 1000
	var items []upapi.CloudStatusService
	for page := int64(1); page <= maxPages; page++ {
		api, err := d.p.api.Checks().ListCloudStatusServices(ctx, upapi.CloudStatusServiceListOptions{
			Page:     page,
			PageSize: pageSize,
			Group:    config.Group.ValueString(),
			Search:   config.Search.ValueString(),
		})
		if err != nil {
			rs.Diagnostics.AddError(
				"uptime_cloudstatus_services read failed",
				fmt.Sprintf("page=%d group=%q search=%q: %s",
					page, config.Group.ValueString(), config.Search.ValueString(), err),
			)
			return
		}
		items = append(items, api.Items...)
		if int64(len(api.Items)) < pageSize || int64(len(items)) >= api.TotalCount {
			break
		}
		if page == maxPages {
			rs.Diagnostics.AddError(
				"uptime_cloudstatus_services read aborted",
				fmt.Sprintf("paginated past %d pages without reaching reported total %d - server may be returning inconsistent counts", maxPages, api.TotalCount),
			)
			return
		}
	}

	model := CloudStatusServicesDataSourceModel{
		ID:       types.StringValue(""),
		Group:    config.Group,
		Search:   config.Search,
		Services: make([]CloudStatusServicesDataSourceItemModel, len(items)),
	}
	for i := range items {
		model.Services[i] = CloudStatusServicesDataSourceItemModel{
			ID:       types.Int64Value(items[i].ID),
			Name:     types.StringValue(items[i].Name),
			Title:    types.StringValue(items[i].Title),
			SubTitle: types.StringValue(items[i].SubTitle),
			GroupID:  types.Int64Value(items[i].GroupID),
			Group:    types.StringValue(items[i].Group),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
