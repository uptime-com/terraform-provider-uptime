package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationTwitterResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationTwitterResourceModel, upapi.IntegrationTwitter, upapi.Integration]{
		api: IntegrationTwitterResourceAPI{provider: p},
		mod: IntegrationTwitterResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_twitter",
			Schema: schema.Schema{
				Description: "Twitter integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"oauth_token": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "OAuth token",
					},
					"oauth_token_secret": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "OAuth token secret",
					},
				},
			},
		},
	}
}

type IntegrationTwitterResourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	URL              types.String `tfsdk:"url"`
	Name             types.String `tfsdk:"name"`
	ContactGroups    types.Set    `tfsdk:"contact_groups"`
	OauthToken       types.String `tfsdk:"oauth_token"`
	OauthTokenSecret types.String `tfsdk:"oauth_token_secret"`
}

func (m IntegrationTwitterResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationTwitterResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationTwitterResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationTwitterResourceModel, diag.Diagnostics) {
	model := *new(IntegrationTwitterResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationTwitterResourceModelAdapter) ToAPIArgument(model IntegrationTwitterResourceModel) (*upapi.IntegrationTwitter, error) {
	return &upapi.IntegrationTwitter{
		Name:             model.Name.ValueString(),
		ContactGroups:    a.ContactGroups(model.ContactGroups),
		OauthToken:       model.OauthToken.ValueString(),
		OauthTokenSecret: model.OauthTokenSecret.ValueString(),
	}, nil
}

func (a IntegrationTwitterResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationTwitterResourceModel, error) {
	return &IntegrationTwitterResourceModel{
		ID:               types.Int64Value(api.PK),
		URL:              types.StringValue(api.URL),
		Name:             types.StringValue(api.Name),
		ContactGroups:    a.ContactGroupsValue(api.ContactGroups),
		OauthToken:       types.StringValue(""),
		OauthTokenSecret: types.StringValue(""),
	}, nil
}

type IntegrationTwitterResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationTwitterResourceAPI) Create(ctx context.Context, arg upapi.IntegrationTwitter) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateTwitter(ctx, arg)
}

func (a IntegrationTwitterResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationTwitterResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationTwitter) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateTwitter(ctx, pk, arg)
}

func (a IntegrationTwitterResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
