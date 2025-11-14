package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewUsersDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return UsersDataSource{p: p}
}

// UsersDataSchema defines the schema for the users data source.
var UsersDataSchema = schema.Schema{
	Description: "Retrieve a list of all users in your Uptime.com account. Users represent team members with access to the account.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"users": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all users in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the user",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the user resource",
					},
					"first_name": schema.StringAttribute{
						Computed:    true,
						Description: "User's first name",
					},
					"last_name": schema.StringAttribute{
						Computed:    true,
						Description: "User's last name",
					},
					"email": schema.StringAttribute{
						Computed:    true,
						Description: "User's email address",
					},
					"is_active": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the user account is active",
					},
					"is_primary": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this is the primary account owner",
					},
					"access_level": schema.StringAttribute{
						Computed:    true,
						Description: "User's access level (e.g., 'admin', 'user')",
					},
					"is_api_enabled": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether API access is enabled for this user",
					},
					"notify_paid_invoices": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether user receives notifications about paid invoices",
					},
					"assigned_subaccounts": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of subaccount URLs assigned to this user",
					},
					"require_two_factor": schema.StringAttribute{
						Computed:    true,
						Description: "Two-factor authentication requirement setting",
					},
					"must_two_factor": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether two-factor authentication is mandatory",
					},
					"timezone": schema.StringAttribute{
						Computed:    true,
						Description: "User's timezone setting",
					},
				},
			},
		},
	},
}

type UsersDataSourceModel struct {
	ID    types.String               `tfsdk:"id"`
	Users []UsersDataSourceItemModel `tfsdk:"users"`
}

type UsersDataSourceItemModel struct {
	ID                  types.Int64  `tfsdk:"id"`
	URL                 types.String `tfsdk:"url"`
	FirstName           types.String `tfsdk:"first_name"`
	LastName            types.String `tfsdk:"last_name"`
	Email               types.String `tfsdk:"email"`
	IsActive            types.Bool   `tfsdk:"is_active"`
	IsPrimary           types.Bool   `tfsdk:"is_primary"`
	AccessLevel         types.String `tfsdk:"access_level"`
	IsAPIEnabled        types.Bool   `tfsdk:"is_api_enabled"`
	NotifyPaidInvoices  types.Bool   `tfsdk:"notify_paid_invoices"`
	AssignedSubaccounts types.List   `tfsdk:"assigned_subaccounts"`
	RequireTwoFactor    types.String `tfsdk:"require_two_factor"`
	MustTwoFactor       types.Bool   `tfsdk:"must_two_factor"`
	Timezone            types.String `tfsdk:"timezone"`
}

var _ datasource.DataSource = &UsersDataSource{}

type UsersDataSource struct {
	p *providerImpl
}

func (d UsersDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_users"
}

func (d UsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = UsersDataSchema
}

func (d UsersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.Users().List(ctx, upapi.UserListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := UsersDataSourceModel{
		ID:    types.StringValue(""),
		Users: make([]UsersDataSourceItemModel, len(api)),
	}

	for i := range api {
		// Convert AssignedSubaccounts slice to types.List
		subaccountsValues := make([]attr.Value, len(api[i].AssignedSubaccounts))
		for j, v := range api[i].AssignedSubaccounts {
			subaccountsValues[j] = types.StringValue(v)
		}
		subaccounts := types.ListNull(types.StringType)
		if len(subaccountsValues) > 0 {
			subaccounts = types.ListValueMust(types.StringType, subaccountsValues)
		}

		model.Users[i] = UsersDataSourceItemModel{
			ID:                  types.Int64Value(api[i].PK),
			URL:                 types.StringValue(api[i].URL),
			FirstName:           types.StringValue(api[i].FirstName),
			LastName:            types.StringValue(api[i].LastName),
			Email:               types.StringValue(api[i].Email),
			IsActive:            types.BoolValue(api[i].IsActive),
			IsPrimary:           types.BoolValue(api[i].IsPrimary),
			AccessLevel:         types.StringValue(api[i].AccessLevel),
			IsAPIEnabled:        types.BoolValue(api[i].IsAPIEnabled),
			NotifyPaidInvoices:  types.BoolValue(api[i].NotifyPaidInvoices),
			AssignedSubaccounts: subaccounts,
			RequireTwoFactor:    types.StringValue(api[i].RequireTwoFactor),
			MustTwoFactor:       types.BoolValue(api[i].MustTwoFactor),
			Timezone:            types.StringValue(api[i].Timezone),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
