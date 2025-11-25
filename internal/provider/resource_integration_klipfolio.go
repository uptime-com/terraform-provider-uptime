package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationKlipfolioResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationKlipfolioResourceModel, upapi.IntegrationKlipfolio, upapi.Integration]{
		api: IntegrationKlipfolioResourceAPI{provider: p},
		mod: IntegrationKlipfolioResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_klipfolio",
			Schema: schema.Schema{
				Description: "Klipfolio integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"api_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Klipfolio API key",
					},
					"data_source_name": schema.StringAttribute{
						Required:    true,
						Description: "Name of the data source",
					},
				},
			},
		},
	}
}

type IntegrationKlipfolioResourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	URL            types.String `tfsdk:"url"`
	Name           types.String `tfsdk:"name"`
	ContactGroups  types.Set    `tfsdk:"contact_groups"`
	APIKey         types.String `tfsdk:"api_key"`
	DataSourceName types.String `tfsdk:"data_source_name"`
}

func (m IntegrationKlipfolioResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationKlipfolioResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationKlipfolioResourceModelAdapter) PreservePlanValues(result *IntegrationKlipfolioResourceModel, plan *IntegrationKlipfolioResourceModel) *IntegrationKlipfolioResourceModel {
	if result.APIKey.IsNull() {
		result.APIKey = plan.APIKey
	}
	if result.DataSourceName.IsNull() {
		result.DataSourceName = plan.DataSourceName
	}
	return result
}

func (a IntegrationKlipfolioResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationKlipfolioResourceModel, diag.Diagnostics) {
	model := *new(IntegrationKlipfolioResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationKlipfolioResourceModelAdapter) ToAPIArgument(model IntegrationKlipfolioResourceModel) (*upapi.IntegrationKlipfolio, error) {
	return &upapi.IntegrationKlipfolio{
		Name:           model.Name.ValueString(),
		ContactGroups:  a.ContactGroups(model.ContactGroups),
		APIKey:         model.APIKey.ValueString(),
		DataSourceName: model.DataSourceName.ValueString(),
	}, nil
}

func (a IntegrationKlipfolioResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationKlipfolioResourceModel, error) {
	return &IntegrationKlipfolioResourceModel{
		ID:             types.Int64Value(api.PK),
		URL:            types.StringValue(api.URL),
		Name:           types.StringValue(api.Name),
		ContactGroups:  a.ContactGroupsValue(api.ContactGroups),
		APIKey:         types.StringNull(),
		DataSourceName: types.StringNull(),
	}, nil
}

type IntegrationKlipfolioResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationKlipfolioResourceAPI) Create(ctx context.Context, arg upapi.IntegrationKlipfolio) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateKlipfolio(ctx, arg)
}

func (a IntegrationKlipfolioResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationKlipfolioResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationKlipfolio) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateKlipfolio(ctx, pk, arg)
}

func (a IntegrationKlipfolioResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
