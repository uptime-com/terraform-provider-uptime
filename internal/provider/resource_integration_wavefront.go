package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationWavefrontResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationWavefrontResourceModel, upapi.IntegrationWavefront, upapi.Integration]{
		api: IntegrationWavefrontResourceAPI{provider: p},
		mod: IntegrationWavefrontResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_wavefront",
			Schema: schema.Schema{
				Description: "Wavefront integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"wavefront_url": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Wavefront instance URL",
					},
					"api_token": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Wavefront API token",
					},
				},
			},
		},
	}
}

type IntegrationWavefrontResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	WavefrontURL  types.String `tfsdk:"wavefront_url"`
	APIToken      types.String `tfsdk:"api_token"`
}

func (m IntegrationWavefrontResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationWavefrontResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationWavefrontResourceModelAdapter) PreservePlanValues(result *IntegrationWavefrontResourceModel, plan *IntegrationWavefrontResourceModel) *IntegrationWavefrontResourceModel {
	if result.WavefrontURL.IsNull() {
		result.WavefrontURL = plan.WavefrontURL
	}
	if result.APIToken.IsNull() {
		result.APIToken = plan.APIToken
	}
	return result
}

func (a IntegrationWavefrontResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationWavefrontResourceModel, diag.Diagnostics) {
	model := *new(IntegrationWavefrontResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationWavefrontResourceModelAdapter) ToAPIArgument(model IntegrationWavefrontResourceModel) (*upapi.IntegrationWavefront, error) {
	return &upapi.IntegrationWavefront{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		WavefrontUrl:  model.WavefrontURL.ValueString(),
		APIToken:      model.APIToken.ValueString(),
	}, nil
}

func (a IntegrationWavefrontResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationWavefrontResourceModel, error) {
	return &IntegrationWavefrontResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		WavefrontURL:  types.StringNull(),
		APIToken:      types.StringNull(),
	}, nil
}

type IntegrationWavefrontResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationWavefrontResourceAPI) Create(ctx context.Context, arg upapi.IntegrationWavefront) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateWavefront(ctx, arg)
}

func (a IntegrationWavefrontResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationWavefrontResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationWavefront) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateWavefront(ctx, pk, arg)
}

func (a IntegrationWavefrontResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
