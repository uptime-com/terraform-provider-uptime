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

func NewIntegrationVictoropsResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationVictoropsResourceModel, upapi.IntegrationVictorops, upapi.Integration]{
		api: IntegrationVictoropsResourceAPI{provider: p},
		mod: IntegrationVictoropsResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_victorops",
			Schema: schema.Schema{
				Description: "VictorOps integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"service_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "VictorOps service API key",
					},
					"routing_key": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "VictorOps routing key",
					},
				},
			},
		},
	}
}

type IntegrationVictoropsResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	ServiceKey    types.String `tfsdk:"service_key"`
	RoutingKey    types.String `tfsdk:"routing_key"`
}

func (m IntegrationVictoropsResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationVictoropsResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationVictoropsResourceModelAdapter) PreservePlanValues(result *IntegrationVictoropsResourceModel, plan *IntegrationVictoropsResourceModel) *IntegrationVictoropsResourceModel {
	if result.ServiceKey.IsNull() {
		result.ServiceKey = plan.ServiceKey
	}
	if result.RoutingKey.IsNull() {
		result.RoutingKey = plan.RoutingKey
	}
	return result
}

func (a IntegrationVictoropsResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationVictoropsResourceModel, diag.Diagnostics) {
	model := *new(IntegrationVictoropsResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationVictoropsResourceModelAdapter) ToAPIArgument(model IntegrationVictoropsResourceModel) (*upapi.IntegrationVictorops, error) {
	return &upapi.IntegrationVictorops{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroups(model.ContactGroups),
		ServiceKey:    model.ServiceKey.ValueString(),
		RoutingKey:    model.RoutingKey.ValueString(),
	}, nil
}

func (a IntegrationVictoropsResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationVictoropsResourceModel, error) {
	return &IntegrationVictoropsResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsValue(api.ContactGroups),
		ServiceKey:    types.StringNull(),
		RoutingKey:    types.StringNull(),
	}, nil
}

type IntegrationVictoropsResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationVictoropsResourceAPI) Create(ctx context.Context, arg upapi.IntegrationVictorops) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateVictorops(ctx, arg)
}

func (a IntegrationVictoropsResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationVictoropsResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationVictorops) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateVictorops(ctx, pk, arg)
}

func (a IntegrationVictoropsResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
