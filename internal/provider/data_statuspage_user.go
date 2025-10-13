package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageUserDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageUserDataSource{p: p}
}

var StatusPageUserDataSchema = schema.Schema{
	Description: "Retrieve a list of all users for a specific status page.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspage_id": schema.Int64Attribute{
			Required:    true,
			Description: "ID of the status page to retrieve users for",
		},
		"users": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all users for the status page",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the user",
					},
					"email": schema.StringAttribute{
						Computed:    true,
						Description: "Email address of the user",
					},
					"first_name": schema.StringAttribute{
						Computed:    true,
						Description: "First name of the user",
					},
					"last_name": schema.StringAttribute{
						Computed:    true,
						Description: "Last name of the user",
					},
					"is_active": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the user is active",
					},
				},
			},
		},
	},
}

type StatusPageUserDataSourceModel struct {
	ID           types.String                        `tfsdk:"id"`
	StatusPageID types.Int64                         `tfsdk:"statuspage_id"`
	Users        []StatusPageUserDataSourceItemModel `tfsdk:"users"`
}

type StatusPageUserDataSourceItemModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	IsActive  types.Bool   `tfsdk:"is_active"`
}

var _ datasource.DataSource = &StatusPageUserDataSource{}

type StatusPageUserDataSource struct {
	p *providerImpl
}

func (d StatusPageUserDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspage_users"
}

func (d StatusPageUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageUserDataSchema
}

func (d StatusPageUserDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config StatusPageUserDataSourceModel
	diags := rq.Config.Get(ctx, &config)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	pk := upapi.PrimaryKey(config.StatusPageID.ValueInt64())
	api, err := d.p.api.StatusPages().Users(pk).List(ctx, upapi.StatusPageUserListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageUserDataSourceModel{
		ID:           types.StringValue(""),
		StatusPageID: config.StatusPageID,
		Users:        make([]StatusPageUserDataSourceItemModel, len(api)),
	}

	for i := range api {
		model.Users[i] = StatusPageUserDataSourceItemModel{
			ID:        types.Int64Value(api[i].PK),
			Email:     types.StringValue(api[i].Email),
			FirstName: types.StringValue(api[i].FirstName),
			LastName:  types.StringValue(api[i].LastName),
			IsActive:  types.BoolValue(api[i].IsActive),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
