package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationGeckoboardResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationGeckoboardResourceModel, upapi.IntegrationGeckoboard, upapi.Integration]{
		api: IntegrationGeckoboardResourceAPI{provider: p},
		mod: IntegrationGeckoboardResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_geckoboard",
			Schema: schema.Schema{
				Description: "Geckoboard integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"api_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Geckoboard API key",
					},
					"dataset_name": schema.StringAttribute{
						Required:    true,
						Description: "Name of the dataset to send data to",
					},
				},
			},
		},
	}
}

type IntegrationGeckoboardResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	APIKey        types.String `tfsdk:"api_key"`
	DatasetName   types.String `tfsdk:"dataset_name"`
}

func (m IntegrationGeckoboardResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationGeckoboardResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationGeckoboardResourceModelAdapter) PreservePlanValues(result *IntegrationGeckoboardResourceModel, plan *IntegrationGeckoboardResourceModel) *IntegrationGeckoboardResourceModel {
	if result.APIKey.IsNull() {
		result.APIKey = plan.APIKey
	}
	if result.DatasetName.IsNull() {
		result.DatasetName = plan.DatasetName
	}
	return result
}

func (a IntegrationGeckoboardResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationGeckoboardResourceModel, diag.Diagnostics) {
	model := *new(IntegrationGeckoboardResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationGeckoboardResourceModelAdapter) ToAPIArgument(model IntegrationGeckoboardResourceModel) (*upapi.IntegrationGeckoboard, error) {
	return &upapi.IntegrationGeckoboard{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		APIKey:        model.APIKey.ValueString(),
		DatasetName:   model.DatasetName.ValueString(),
	}, nil
}

func (a IntegrationGeckoboardResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationGeckoboardResourceModel, error) {
	return &IntegrationGeckoboardResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		APIKey:        types.StringNull(),
		DatasetName:   types.StringNull(),
	}, nil
}

type IntegrationGeckoboardResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationGeckoboardResourceAPI) Create(ctx context.Context, arg upapi.IntegrationGeckoboard) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateGeckoboard(ctx, arg)
}

func (a IntegrationGeckoboardResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationGeckoboardResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationGeckoboard) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateGeckoboard(ctx, pk, arg)
}

func (a IntegrationGeckoboardResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
