package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationStatuspageResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationStatuspageResourceModel, upapi.IntegrationStatuspage, upapi.Integration]{
		api: IntegrationStatuspageResourceAPI{provider: p},
		mod: IntegrationStatuspageResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_statuspage",
			Schema: schema.Schema{
				Description: "Statuspage.io integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"api_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Statuspage.io API key",
					},
					"page": schema.StringAttribute{
						Required:    true,
						Description: "Statuspage.io page ID",
					},
					"component": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Component ID to update",
					},
					"metric": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Metric ID to update",
					},
				},
			},
		},
	}
}

type IntegrationStatuspageResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	APIKey        types.String `tfsdk:"api_key"`
	Page          types.String `tfsdk:"page"`
	Component     types.String `tfsdk:"component"`
	Metric        types.String `tfsdk:"metric"`
}

func (m IntegrationStatuspageResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationStatuspageResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationStatuspageResourceModelAdapter) PreservePlanValues(result *IntegrationStatuspageResourceModel, plan *IntegrationStatuspageResourceModel) *IntegrationStatuspageResourceModel {
	if result.APIKey.IsNull() {
		result.APIKey = plan.APIKey
	}
	if result.Page.IsNull() {
		result.Page = plan.Page
	}
	if result.Component.IsNull() {
		result.Component = plan.Component
	}
	if result.Metric.IsNull() {
		result.Metric = plan.Metric
	}
	return result
}

func (a IntegrationStatuspageResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationStatuspageResourceModel, diag.Diagnostics) {
	model := *new(IntegrationStatuspageResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationStatuspageResourceModelAdapter) ToAPIArgument(model IntegrationStatuspageResourceModel) (*upapi.IntegrationStatuspage, error) {
	return &upapi.IntegrationStatuspage{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroups(model.ContactGroups),
		APIKey:        model.APIKey.ValueString(),
		Page:          model.Page.ValueString(),
		Component:     model.Component.ValueString(),
		Metric:        model.Metric.ValueString(),
	}, nil
}

func (a IntegrationStatuspageResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationStatuspageResourceModel, error) {
	return &IntegrationStatuspageResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsValue(api.ContactGroups),
		APIKey:        types.StringNull(),
		Page:          types.StringNull(),
		Component:     types.StringNull(),
		Metric:        types.StringNull(),
	}, nil
}

type IntegrationStatuspageResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationStatuspageResourceAPI) Create(ctx context.Context, arg upapi.IntegrationStatuspage) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateStatuspage(ctx, arg)
}

func (a IntegrationStatuspageResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationStatuspageResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationStatuspage) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateStatuspage(ctx, pk, arg)
}

func (a IntegrationStatuspageResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
