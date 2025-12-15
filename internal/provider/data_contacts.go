package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewContactsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return ContactsDataSource{p: p}
}

// ContactsDataSchema defines the schema for the contacts data source.
var ContactsDataSchema = schema.Schema{
	Description: "Retrieve a list of all contacts configured in your Uptime.com account. Contacts define notification recipients for alerts.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"contacts": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all contacts in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the contact",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the contact resource",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Human-readable name for the contact",
					},
					"sms_list": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of phone numbers for SMS notifications",
					},
					"email_list": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of email addresses for email notifications",
					},
					"phonecall_list": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of phone numbers for voice call notifications",
					},
					"integrations": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of integration URLs for third-party notifications",
					},
					"push_notification_profiles": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of push notification profile URLs for mobile notifications",
					},
				},
			},
		},
	},
}

type ContactsDataSourceModel struct {
	ID       types.String                  `tfsdk:"id"`
	Contacts []ContactsDataSourceItemModel `tfsdk:"contacts"`
}

type ContactsDataSourceItemModel struct {
	ID                       types.Int64  `tfsdk:"id"`
	URL                      types.String `tfsdk:"url"`
	Name                     types.String `tfsdk:"name"`
	SMSList                  types.List   `tfsdk:"sms_list"`
	EmailList                types.List   `tfsdk:"email_list"`
	PhonecallList            types.List   `tfsdk:"phonecall_list"`
	Integrations             types.List   `tfsdk:"integrations"`
	PushNotificationProfiles types.List   `tfsdk:"push_notification_profiles"`
}

var _ datasource.DataSource = &ContactsDataSource{}

type ContactsDataSource struct {
	p *providerImpl
}

func (d ContactsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_contacts"
}

func (d ContactsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = ContactsDataSchema
}

func (d ContactsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.Contacts().List(ctx, upapi.ContactListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := ContactsDataSourceModel{
		ID:       types.StringValue(""),
		Contacts: make([]ContactsDataSourceItemModel, len(api.Items)),
	}

	for i := range api.Items {
		// Convert string slices to types.List
		smsListValues := make([]attr.Value, len(api.Items[i].SmsList))
		for j, v := range api.Items[i].SmsList {
			smsListValues[j] = types.StringValue(v)
		}
		smsList := types.ListNull(types.StringType)
		if len(smsListValues) > 0 {
			smsList = types.ListValueMust(types.StringType, smsListValues)
		}

		emailListValues := make([]attr.Value, len(api.Items[i].EmailList))
		for j, v := range api.Items[i].EmailList {
			emailListValues[j] = types.StringValue(v)
		}
		emailList := types.ListNull(types.StringType)
		if len(emailListValues) > 0 {
			emailList = types.ListValueMust(types.StringType, emailListValues)
		}

		phonecallListValues := make([]attr.Value, len(api.Items[i].PhonecallList))
		for j, v := range api.Items[i].PhonecallList {
			phonecallListValues[j] = types.StringValue(v)
		}
		phonecallList := types.ListNull(types.StringType)
		if len(phonecallListValues) > 0 {
			phonecallList = types.ListValueMust(types.StringType, phonecallListValues)
		}

		integrationsValues := make([]attr.Value, len(api.Items[i].Integrations))
		for j, v := range api.Items[i].Integrations {
			integrationsValues[j] = types.StringValue(v)
		}
		integrations := types.ListNull(types.StringType)
		if len(integrationsValues) > 0 {
			integrations = types.ListValueMust(types.StringType, integrationsValues)
		}

		pushProfilesValues := make([]attr.Value, len(api.Items[i].PushNotificationProfiles))
		for j, v := range api.Items[i].PushNotificationProfiles {
			pushProfilesValues[j] = types.StringValue(v)
		}
		pushProfiles := types.ListNull(types.StringType)
		if len(pushProfilesValues) > 0 {
			pushProfiles = types.ListValueMust(types.StringType, pushProfilesValues)
		}

		model.Contacts[i] = ContactsDataSourceItemModel{
			ID:                       types.Int64Value(api.Items[i].PK),
			URL:                      types.StringValue(api.Items[i].URL),
			Name:                     types.StringValue(api.Items[i].Name),
			SMSList:                  smsList,
			EmailList:                emailList,
			PhonecallList:            phonecallList,
			Integrations:             integrations,
			PushNotificationProfiles: pushProfiles,
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
