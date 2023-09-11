package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewContactResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[contactResourceModel, upapi.Contact, upapi.Contact]{
		api: &contactResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "contact",
			Schema:         contactResourceSchema,
		},
	}
}

var contactResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id":   IDAttribute(),
		"url":  URLAttribute(),
		"name": NameAttribute(),
		"sms_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     StringSetEmptyDefault(),
		},
		"email_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     StringSetEmptyDefault(),
		},
		"phonecall_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     StringSetEmptyDefault(),
		},
		"integrations": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     StringSetEmptyDefault(),
		},
		"push_notification_profiles": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Optional:    true,
			Default:     StringSetEmptyDefault(),
		},
	},
}

type contactResourceModel struct {
	ID                       types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                      types.String `tfsdk:"url"  ref:"URL,opt"`
	Name                     types.String `tfsdk:"name"`
	SmsList                  types.Set    `tfsdk:"sms_list"`
	EmailList                types.Set    `tfsdk:"email_list"`
	PhonecallList            types.Set    `tfsdk:"phonecall_list"`
	Integrations             types.Set    `tfsdk:"integrations"`
	PushNotificationProfiles types.Set    `tfsdk:"push_notification_profiles"`
}

var _ genericResourceAPI[upapi.Contact, upapi.Contact] = (*contactResourceAPI)(nil)

type contactResourceAPI struct {
	provider *providerImpl
}

func (c *contactResourceAPI) Create(ctx context.Context, arg upapi.Contact) (*upapi.Contact, error) {
	obj, err := c.provider.api.Contacts().Create(ctx, arg)
	return obj, err
}

func (c *contactResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Contact, error) {
	obj, err := c.provider.api.Contacts().Get(ctx, pk)
	return obj, err
}

func (c *contactResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.Contact) (*upapi.Contact, error) {
	obj, err := c.provider.api.Contacts().Update(ctx, pk, arg)
	return obj, err
}

func (c *contactResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Contacts().Delete(ctx, pk)
}
