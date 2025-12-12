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

func NewIntegrationCachetResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationCachetResourceModel, upapi.IntegrationCachet, upapi.Integration](
		IntegrationCachetResourceAPI{provider: p},
		IntegrationCachetResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_cachet",
			Schema: schema.Schema{
				Description: "Cachet integration resource. Import using the integration ID: `terraform import uptime_integration_cachet.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"cachet_url": schema.StringAttribute{
						Required:    true,
						Description: "The URL of your Cachet instance",
					},
					"token": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Cachet API token",
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
		ImportStateSimpleID,
	)
}

type IntegrationCachetResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	CachetURL     types.String `tfsdk:"cachet_url"`
	Token         types.String `tfsdk:"token"`
	Component     types.String `tfsdk:"component"`
	Metric        types.String `tfsdk:"metric"`
}

func (m IntegrationCachetResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationCachetResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationCachetResourceModelAdapter) PreservePlanValues(result *IntegrationCachetResourceModel, plan *IntegrationCachetResourceModel) *IntegrationCachetResourceModel {
	if result.CachetURL.IsNull() {
		result.CachetURL = plan.CachetURL
	}
	if result.Token.IsNull() {
		result.Token = plan.Token
	}
	if result.Component.IsNull() {
		result.Component = plan.Component
	}
	if result.Metric.IsNull() {
		result.Metric = plan.Metric
	}
	return result
}

func (a IntegrationCachetResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationCachetResourceModel, diag.Diagnostics) {
	model := *new(IntegrationCachetResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationCachetResourceModelAdapter) ToAPIArgument(model IntegrationCachetResourceModel) (*upapi.IntegrationCachet, error) {
	return &upapi.IntegrationCachet{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		CachetURL:     model.CachetURL.ValueString(),
		Token:         model.Token.ValueString(),
		Component:     model.Component.ValueString(),
		Metric:        model.Metric.ValueString(),
	}, nil
}

func (a IntegrationCachetResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationCachetResourceModel, error) {
	return &IntegrationCachetResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		CachetURL:     types.StringNull(),
		Token:         types.StringNull(),
		Component:     types.StringNull(),
		Metric:        types.StringNull(),
	}, nil
}

type IntegrationCachetResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationCachetResourceAPI) Create(ctx context.Context, arg upapi.IntegrationCachet) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateCachet(ctx, arg)
}

func (a IntegrationCachetResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationCachetResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationCachet) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateCachet(ctx, pk, arg)
}

func (a IntegrationCachetResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
