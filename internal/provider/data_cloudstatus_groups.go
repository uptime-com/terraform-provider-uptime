package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCloudStatusGroupsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return CloudStatusGroupsDataSource{p: p}
}

var CloudStatusGroupsDataSchema = schema.Schema{
	Description: "Look up Cloud Status provider groups - the vendors/providers that can be referenced by a cloudstatus check. Pair with `uptime_cloudstatus_services` to discover the services under each group.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"search": schema.StringAttribute{
			Optional:    true,
			Description: "Case-insensitive substring matched against the provider group name",
		},
		"groups": schema.ListNestedAttribute{
			Computed:    true,
			Description: "Matching Cloud Status provider groups",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the provider group",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the provider group",
					},
				},
			},
		},
	},
}

type CloudStatusGroupsDataSourceModel struct {
	ID     types.String                           `tfsdk:"id"`
	Search types.String                           `tfsdk:"search"`
	Groups []CloudStatusGroupsDataSourceItemModel `tfsdk:"groups"`
}

type CloudStatusGroupsDataSourceItemModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

var _ datasource.DataSource = &CloudStatusGroupsDataSource{}

type CloudStatusGroupsDataSource struct {
	p *providerImpl
}

func (d CloudStatusGroupsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_cloudstatus_groups"
}

func (d CloudStatusGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = CloudStatusGroupsDataSchema
}

func (d CloudStatusGroupsDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config CloudStatusGroupsDataSourceModel
	rs.Diagnostics.Append(rq.Config.Get(ctx, &config)...)
	if rs.Diagnostics.HasError() {
		return
	}

	const pageSize int64 = 100
	const maxPages int64 = 1000
	var items []upapi.CloudStatusGroupListItem
	for page := int64(1); page <= maxPages; page++ {
		api, err := d.p.api.Checks().ListCloudStatusGroups(ctx, upapi.CloudStatusGroupListOptions{
			Page:     page,
			PageSize: pageSize,
			Search:   config.Search.ValueString(),
		})
		if err != nil {
			rs.Diagnostics.AddError(
				"uptime_cloudstatus_groups read failed",
				fmt.Sprintf("page=%d search=%q: %s", page, config.Search.ValueString(), err),
			)
			return
		}
		items = append(items, api.Items...)
		if int64(len(api.Items)) < pageSize || int64(len(items)) >= api.TotalCount {
			break
		}
		if page == maxPages {
			rs.Diagnostics.AddError(
				"uptime_cloudstatus_groups read aborted",
				fmt.Sprintf("paginated past %d pages without reaching reported total %d - server may be returning inconsistent counts", maxPages, api.TotalCount),
			)
			return
		}
	}

	model := CloudStatusGroupsDataSourceModel{
		ID:     types.StringValue(""),
		Search: config.Search,
		Groups: make([]CloudStatusGroupsDataSourceItemModel, len(items)),
	}
	for i := range items {
		model.Groups[i] = CloudStatusGroupsDataSourceItemModel{
			ID:   types.Int64Value(items[i].ID),
			Name: types.StringValue(items[i].Name),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
