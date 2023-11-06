package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewContactResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[ContactResourceModel, upapi.Contact, upapi.Contact]{
		api: ContactResourceAPI{provider: p},
		mod: ContactResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "contact",
			Schema:         contactResourceSchema,
		},
	}
}

var contactResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id":   IDSchemaAttribute(),
		"url":  URLSchemaAttribute(),
		"name": NameSchemaAttribute(),
		"sms_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
		},
		"email_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
		},
		"phonecall_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
		},
		"integrations": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
		},
		"push_notification_profiles": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
		},
	},
}

type ContactResourceModel struct {
	ID                       types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                      types.String `tfsdk:"url"  ref:"URL,opt"`
	Name                     types.String `tfsdk:"name"`
	SMSList                  types.Set    `tfsdk:"sms_list"`
	EmailList                types.Set    `tfsdk:"email_list"`
	PhonecallList            types.Set    `tfsdk:"phonecall_list"`
	Integrations             types.Set    `tfsdk:"integrations"`
	PushNotificationProfiles types.Set    `tfsdk:"push_notification_profiles"`
}

func (m ContactResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type ContactResourceModelAdapter struct {
	SetAttributeAdapter[string]
}

func (a ContactResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*ContactResourceModel, diag.Diagnostics) {
	model := *new(ContactResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a ContactResourceModelAdapter) ToAPIArgument(model ContactResourceModel) (*upapi.Contact, error) {
	api := upapi.Contact{
		Name:                     model.Name.ValueString(),
		SmsList:                  a.Slice(model.SMSList),
		EmailList:                a.Slice(model.EmailList),
		PhonecallList:            a.Slice(model.PhonecallList),
		Integrations:             a.Slice(model.Integrations),
		PushNotificationProfiles: a.Slice(model.PushNotificationProfiles),
	}
	return &api, nil
}

func (a ContactResourceModelAdapter) FromAPIResult(api upapi.Contact) (*ContactResourceModel, error) {
	model := ContactResourceModel{
		ID:                       types.Int64Value(api.PK),
		URL:                      types.StringValue(api.URL),
		Name:                     types.StringValue(api.Name),
		SMSList:                  a.SliceValue(api.SmsList),
		EmailList:                a.SliceValue(api.EmailList),
		PhonecallList:            a.SliceValue(api.PhonecallList),
		Integrations:             a.SliceValue(api.Integrations),
		PushNotificationProfiles: a.SliceValue(api.PushNotificationProfiles),
	}
	return &model, nil
}

type ContactResourceAPI struct {
	provider *providerImpl
}

func (c ContactResourceAPI) Create(ctx context.Context, arg upapi.Contact) (*upapi.Contact, error) {
	obj, err := c.provider.api.Contacts().Create(ctx, arg)
	return obj, err
}

func (c ContactResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Contact, error) {
	obj, err := c.provider.api.Contacts().Get(ctx, pk)
	return obj, err
}

func (c ContactResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.Contact) (*upapi.Contact, error) {
	obj, err := c.provider.api.Contacts().Update(ctx, pk, arg)
	return obj, err
}

func (c ContactResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Contacts().Delete(ctx, pk)
}
