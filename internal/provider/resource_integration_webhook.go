package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationWebhookResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationWebhookResourceModel, upapi.IntegrationWebhook, upapi.Integration]{
		api: IntegrationWebhookResourceAPI{provider: p},
		mod: IntegrationWebhookResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_webhook",
			Schema: schema.Schema{
				Description: "Webhook integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"postback_url": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The webhook URL to POST to",
					},
					"headers": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Custom headers to send with the webhook request (JSON format)",
					},
					"use_legacy_payload": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Use legacy payload format",
					},
				},
			},
		},
	}
}

type IntegrationWebhookResourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	URL              types.String `tfsdk:"url"`
	Name             types.String `tfsdk:"name"`
	ContactGroups    types.Set    `tfsdk:"contact_groups"`
	PostbackURL      types.String `tfsdk:"postback_url"`
	Headers          types.String `tfsdk:"headers"`
	UseLegacyPayload types.Bool   `tfsdk:"use_legacy_payload"`
}

func (m IntegrationWebhookResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationWebhookResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationWebhookResourceModelAdapter) PreservePlanValues(result *IntegrationWebhookResourceModel, plan *IntegrationWebhookResourceModel) *IntegrationWebhookResourceModel {
	// Preserve postback_url from plan if API didn't return it
	if result.PostbackURL.IsNull() {
		result.PostbackURL = plan.PostbackURL
	}
	// Preserve headers from plan if API didn't return it
	if result.Headers.IsNull() {
		result.Headers = plan.Headers
	}
	// Preserve use_legacy_payload from plan if API didn't return it
	if result.UseLegacyPayload.IsNull() {
		result.UseLegacyPayload = plan.UseLegacyPayload
	}
	return result
}

func (a IntegrationWebhookResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationWebhookResourceModel, diag.Diagnostics) {
	model := *new(IntegrationWebhookResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationWebhookResourceModelAdapter) ToAPIArgument(model IntegrationWebhookResourceModel) (*upapi.IntegrationWebhook, error) {
	return &upapi.IntegrationWebhook{
		Name:             model.Name.ValueString(),
		ContactGroups:    a.ContactGroupsSlice(model.ContactGroups),
		PostbackUrl:      model.PostbackURL.ValueString(),
		Headers:          model.Headers.ValueString(),
		UseLegacyPayload: model.UseLegacyPayload.ValueBool(),
	}, nil
}

func (a IntegrationWebhookResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationWebhookResourceModel, error) {
	return &IntegrationWebhookResourceModel{
		ID:               types.Int64Value(api.PK),
		URL:              types.StringValue(api.URL),
		Name:             types.StringValue(api.Name),
		ContactGroups:    a.ContactGroupsSliceValue(api.ContactGroups),
		PostbackURL:      types.StringNull(),
		Headers:          types.StringNull(),
		UseLegacyPayload: types.BoolNull(),
	}, nil
}

type IntegrationWebhookResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationWebhookResourceAPI) Create(ctx context.Context, arg upapi.IntegrationWebhook) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateWebhook(ctx, arg)
}

func (a IntegrationWebhookResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationWebhookResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationWebhook) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateWebhook(ctx, pk, arg)
}

func (a IntegrationWebhookResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
