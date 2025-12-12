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

func NewIntegrationDatadogResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationDatadogResourceModel, upapi.IntegrationDatadog, upapi.Integration](
		IntegrationDatadogResourceAPI{provider: p},
		IntegrationDatadogResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_datadog",
			Schema: schema.Schema{
				Description: "Datadog integration resource. Import using the integration ID: `terraform import uptime_integration_datadog.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"api_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Datadog API key",
					},
					"app_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Datadog application key",
					},
					"region": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Datadog region (e.g., 'us', 'eu')",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type IntegrationDatadogResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	APIKey        types.String `tfsdk:"api_key"`
	APPKey        types.String `tfsdk:"app_key"`
	Region        types.String `tfsdk:"region"`
}

func (m IntegrationDatadogResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationDatadogResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationDatadogResourceModelAdapter) PreservePlanValues(result *IntegrationDatadogResourceModel, plan *IntegrationDatadogResourceModel) *IntegrationDatadogResourceModel {
	if result.APIKey.IsNull() {
		result.APIKey = plan.APIKey
	}
	if result.APPKey.IsNull() {
		result.APPKey = plan.APPKey
	}
	if result.Region.IsNull() {
		result.Region = plan.Region
	}
	return result
}

func (a IntegrationDatadogResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationDatadogResourceModel, diag.Diagnostics) {
	model := *new(IntegrationDatadogResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationDatadogResourceModelAdapter) ToAPIArgument(model IntegrationDatadogResourceModel) (*upapi.IntegrationDatadog, error) {
	return &upapi.IntegrationDatadog{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		APIKey:        model.APIKey.ValueString(),
		APPKey:        model.APPKey.ValueString(),
		Region:        model.Region.ValueString(),
	}, nil
}

func (a IntegrationDatadogResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationDatadogResourceModel, error) {
	return &IntegrationDatadogResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		APIKey:        types.StringNull(),
		APPKey:        types.StringNull(),
		Region:        types.StringNull(),
	}, nil
}

type IntegrationDatadogResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationDatadogResourceAPI) Create(ctx context.Context, arg upapi.IntegrationDatadog) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateDatadog(ctx, arg)
}

func (a IntegrationDatadogResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationDatadogResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationDatadog) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateDatadog(ctx, pk, arg)
}

func (a IntegrationDatadogResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
