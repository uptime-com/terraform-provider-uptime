package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationPushbulletResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationPushbulletResourceModel, upapi.IntegrationPushbullet, upapi.Integration](
		IntegrationPushbulletResourceAPI{provider: p},
		IntegrationPushbulletResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_pushbullet",
			Schema: schema.Schema{
				Description: "Pushbullet integration resource. Import using the integration ID: `terraform import uptime_integration_pushbullet.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"email": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Email address to send notifications to",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type IntegrationPushbulletResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	Email         types.String `tfsdk:"email"`
}

func (m IntegrationPushbulletResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationPushbulletResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationPushbulletResourceModelAdapter) PreservePlanValues(result *IntegrationPushbulletResourceModel, plan *IntegrationPushbulletResourceModel) *IntegrationPushbulletResourceModel {
	if result.Email.IsNull() {
		result.Email = plan.Email
	}
	return result
}

func (a IntegrationPushbulletResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationPushbulletResourceModel, diag.Diagnostics) {
	model := *new(IntegrationPushbulletResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationPushbulletResourceModelAdapter) ToAPIArgument(model IntegrationPushbulletResourceModel) (*upapi.IntegrationPushbullet, error) {
	return &upapi.IntegrationPushbullet{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		Email:         model.Email.ValueString(),
	}, nil
}

func (a IntegrationPushbulletResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationPushbulletResourceModel, error) {
	return &IntegrationPushbulletResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		Email:         types.StringNull(),
	}, nil
}

type IntegrationPushbulletResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationPushbulletResourceAPI) Create(ctx context.Context, arg upapi.IntegrationPushbullet) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreatePushbullet(ctx, arg)
}

func (a IntegrationPushbulletResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationPushbulletResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationPushbullet) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdatePushbullet(ctx, pk, arg)
}

func (a IntegrationPushbulletResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
