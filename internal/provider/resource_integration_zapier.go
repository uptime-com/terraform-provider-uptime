package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationZapierResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationZapierResourceModel, upapi.IntegrationZapier, upapi.Integration]{
		api: IntegrationZapierResourceAPI{provider: p},
		mod: IntegrationZapierResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_zapier",
			Schema: schema.Schema{
				Description: "Zapier integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"webhook_url": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Zapier webhook URL",
					},
				},
			},
		},
	}
}

type IntegrationZapierResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	WebhookURL    types.String `tfsdk:"webhook_url"`
}

func (m IntegrationZapierResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationZapierResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationZapierResourceModelAdapter) PreservePlanValues(result *IntegrationZapierResourceModel, plan *IntegrationZapierResourceModel) *IntegrationZapierResourceModel {
	if result.WebhookURL.IsNull() {
		result.WebhookURL = plan.WebhookURL
	}
	return result
}

func (a IntegrationZapierResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationZapierResourceModel, diag.Diagnostics) {
	model := *new(IntegrationZapierResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationZapierResourceModelAdapter) ToAPIArgument(model IntegrationZapierResourceModel) (*upapi.IntegrationZapier, error) {
	return &upapi.IntegrationZapier{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		WebhookUrl:    model.WebhookURL.ValueString(),
	}, nil
}

func (a IntegrationZapierResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationZapierResourceModel, error) {
	return &IntegrationZapierResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		WebhookURL:    types.StringNull(),
	}, nil
}

type IntegrationZapierResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationZapierResourceAPI) Create(ctx context.Context, arg upapi.IntegrationZapier) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateZapier(ctx, arg)
}

func (a IntegrationZapierResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationZapierResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationZapier) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateZapier(ctx, pk, arg)
}

func (a IntegrationZapierResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
