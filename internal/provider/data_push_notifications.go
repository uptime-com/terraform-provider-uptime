package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewPushNotificationsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return PushNotificationsDataSource{p: p}
}

// PushNotificationsDataSchema defines the schema for the push notifications data source.
var PushNotificationsDataSchema = schema.Schema{
	Description: "Retrieve a list of all push notification profiles configured in your Uptime.com account. Push notification profiles represent mobile devices registered to receive alerts.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"profiles": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all push notification profiles in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the push notification profile",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the push notification profile",
					},
					"uuid": schema.StringAttribute{
						Computed:    true,
						Description: "Unique device identifier",
					},
					"device_name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the mobile device",
					},
					"display_name": schema.StringAttribute{
						Computed:    true,
						Description: "Display name for the push notification profile",
					},
					"user": schema.StringAttribute{
						Computed:    true,
						Description: "User associated with the push notification profile",
					},
					"contact_groups": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of contact groups that trigger notifications to this device",
					},
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the profile was created",
					},
					"modified_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the profile was last modified",
					},
				},
			},
		},
	},
}

type PushNotificationsDataSourceModel struct {
	ID       types.String                              `tfsdk:"id"`
	Profiles []PushNotificationsDataSourceProfileModel `tfsdk:"profiles"`
}

type PushNotificationsDataSourceProfileModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	UUID          types.String `tfsdk:"uuid"`
	DeviceName    types.String `tfsdk:"device_name"`
	DisplayName   types.String `tfsdk:"display_name"`
	User          types.String `tfsdk:"user"`
	ContactGroups types.List   `tfsdk:"contact_groups"`
	CreatedAt     types.String `tfsdk:"created_at"`
	ModifiedAt    types.String `tfsdk:"modified_at"`
}

var _ datasource.DataSource = &PushNotificationsDataSource{}

type PushNotificationsDataSource struct {
	p *providerImpl
}

func (d PushNotificationsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_push_notifications"
}

func (d PushNotificationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = PushNotificationsDataSchema
}

func (d PushNotificationsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.PushNotifications().List(ctx, upapi.PushNotificationProfileListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := PushNotificationsDataSourceModel{
		ID:       types.StringValue(""),
		Profiles: make([]PushNotificationsDataSourceProfileModel, len(api)),
	}

	for i := range api {
		// Convert ContactGroups to types.List
		var contactGroupsList types.List
		if len(api[i].ContactGroups) > 0 {
			elements := make([]attr.Value, len(api[i].ContactGroups))
			for j, cg := range api[i].ContactGroups {
				elements[j] = types.StringValue(cg)
			}
			contactGroupsList = types.ListValueMust(types.StringType, elements)
		} else {
			contactGroupsList = types.ListNull(types.StringType)
		}

		model.Profiles[i] = PushNotificationsDataSourceProfileModel{
			ID:            types.Int64Value(api[i].PK),
			URL:           types.StringValue(api[i].URL),
			UUID:          types.StringValue(api[i].UUID),
			DeviceName:    types.StringValue(api[i].DeviceName),
			DisplayName:   types.StringValue(api[i].DisplayName),
			User:          types.StringValue(api[i].User),
			ContactGroups: contactGroupsList,
			CreatedAt:     types.StringValue(api[i].CreatedAt),
			ModifiedAt:    types.StringValue(api[i].ModifiedAt),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
