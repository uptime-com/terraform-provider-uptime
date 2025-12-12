package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationPagerdutyResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationPagerdutyResourceModel, upapi.IntegrationPagerduty, upapi.Integration](
		IntegrationPagerdutyResourceAPI{provider: p},
		IntegrationPagerdutyResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_pagerduty",
			Schema: schema.Schema{
				Description: "PagerDuty integration resource. Import using the integration ID: `terraform import uptime_integration_pagerduty.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"service_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "PagerDuty service integration key",
					},
					"auto_resolve": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Automatically resolve incident once the check is back up",
						Default:     booldefault.StaticBool(false),
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type IntegrationPagerdutyResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	ServiceKey    types.String `tfsdk:"service_key"`
	AutoResolve   types.Bool   `tfsdk:"auto_resolve"`
}

func (m IntegrationPagerdutyResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationPagerdutyResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationPagerdutyResourceModelAdapter) PreservePlanValues(result *IntegrationPagerdutyResourceModel, plan *IntegrationPagerdutyResourceModel) *IntegrationPagerdutyResourceModel {
	if result.ServiceKey.IsNull() {
		result.ServiceKey = plan.ServiceKey
	}
	return result
}

func (a IntegrationPagerdutyResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationPagerdutyResourceModel, diag.Diagnostics) {
	model := *new(IntegrationPagerdutyResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationPagerdutyResourceModelAdapter) ToAPIArgument(model IntegrationPagerdutyResourceModel) (*upapi.IntegrationPagerduty, error) {
	return &upapi.IntegrationPagerduty{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		ServiceKey:    model.ServiceKey.ValueString(),
		Autoresolve:   model.AutoResolve.ValueBool(),
	}, nil
}

func (a IntegrationPagerdutyResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationPagerdutyResourceModel, error) {
	return &IntegrationPagerdutyResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		ServiceKey:    types.StringNull(),
		AutoResolve:   types.BoolValue(api.Autoresolve),
	}, nil
}

type IntegrationPagerdutyResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationPagerdutyResourceAPI) Create(ctx context.Context, arg upapi.IntegrationPagerduty) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreatePagerduty(ctx, arg)
}

func (a IntegrationPagerdutyResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationPagerdutyResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationPagerduty) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdatePagerduty(ctx, pk, arg)
}

func (a IntegrationPagerdutyResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
