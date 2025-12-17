package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewPushNotificationProfilesDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return PushNotificationProfilesDataSource{p: p}
}

// PushNotificationProfilesDataSchema defines the schema for the push notification profiles data source.
var PushNotificationProfilesDataSchema = schema.Schema{
	Description: "Retrieve a list of all push notification profiles configured in your Uptime.com account. Push notification profiles represent mobile devices registered to receive push notifications.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"push_notification_profiles": schema.ListNestedAttribute{
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
						Description: "API URL for the push notification profile resource",
					},
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the profile was created",
					},
					"modified_at": schema.StringAttribute{
						Computed:    true,
						Description: "Timestamp when the profile was last modified",
					},
					"uuid": schema.StringAttribute{
						Computed:    true,
						Description: "Unique identifier for the device",
					},
					"user": schema.StringAttribute{
						Computed:    true,
						Description: "URL of the associated user",
					},
					"device_name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the mobile device",
					},
					"display_name": schema.StringAttribute{
						Computed:    true,
						Description: "Display name for the push notification profile",
					},
					"contact_groups": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of contact group URLs associated with this profile",
					},
				},
			},
		},
	},
}

type PushNotificationProfilesDataSourceModel struct {
	ID                       types.String                                  `tfsdk:"id"`
	PushNotificationProfiles []PushNotificationProfilesDataSourceItemModel `tfsdk:"push_notification_profiles"`
}

type PushNotificationProfilesDataSourceItemModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	CreatedAt     types.String `tfsdk:"created_at"`
	ModifiedAt    types.String `tfsdk:"modified_at"`
	UUID          types.String `tfsdk:"uuid"`
	User          types.String `tfsdk:"user"`
	DeviceName    types.String `tfsdk:"device_name"`
	DisplayName   types.String `tfsdk:"display_name"`
	ContactGroups types.List   `tfsdk:"contact_groups"`
}

var _ datasource.DataSource = &PushNotificationProfilesDataSource{}

type PushNotificationProfilesDataSource struct {
	p *providerImpl
}

func (d PushNotificationProfilesDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_push_notification_profiles"
}

func (d PushNotificationProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = PushNotificationProfilesDataSchema
}

func (d PushNotificationProfilesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.PushNotifications().List(ctx, upapi.PushNotificationProfileListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := PushNotificationProfilesDataSourceModel{
		ID:                       types.StringValue(""),
		PushNotificationProfiles: make([]PushNotificationProfilesDataSourceItemModel, len(api.Items)),
	}

	for i := range api.Items {
		// Convert ContactGroups slice to types.List
		contactGroupsValues := make([]attr.Value, len(api.Items[i].ContactGroups))
		for j, v := range api.Items[i].ContactGroups {
			contactGroupsValues[j] = types.StringValue(v)
		}
		contactGroups := types.ListNull(types.StringType)
		if len(contactGroupsValues) > 0 {
			contactGroups = types.ListValueMust(types.StringType, contactGroupsValues)
		}

		model.PushNotificationProfiles[i] = PushNotificationProfilesDataSourceItemModel{
			ID:            types.Int64Value(api.Items[i].PK),
			URL:           types.StringValue(api.Items[i].URL),
			CreatedAt:     types.StringValue(api.Items[i].CreatedAt),
			ModifiedAt:    types.StringValue(api.Items[i].ModifiedAt),
			UUID:          types.StringValue(api.Items[i].UUID),
			User:          types.StringValue(api.Items[i].User),
			DeviceName:    types.StringValue(api.Items[i].DeviceName),
			DisplayName:   types.StringValue(api.Items[i].DisplayName),
			ContactGroups: contactGroups,
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
