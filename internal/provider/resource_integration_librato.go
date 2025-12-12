package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationLibratoResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationLibratoResourceModel, upapi.IntegrationLibrato, upapi.Integration](
		IntegrationLibratoResourceAPI{provider: p},
		IntegrationLibratoResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_librato",
			Schema: schema.Schema{
				Description: "Librato integration resource. Import using the integration ID: `terraform import uptime_integration_librato.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"email": schema.StringAttribute{
						Required:    true,
						Description: "Librato account email",
					},
					"api_token": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Librato API token",
					},
					"metric_name": schema.StringAttribute{
						Required:    true,
						Description: "Name of the metric to send data to",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type IntegrationLibratoResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	Email         types.String `tfsdk:"email"`
	APIToken      types.String `tfsdk:"api_token"`
	MetricName    types.String `tfsdk:"metric_name"`
}

func (m IntegrationLibratoResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationLibratoResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationLibratoResourceModelAdapter) PreservePlanValues(result *IntegrationLibratoResourceModel, plan *IntegrationLibratoResourceModel) *IntegrationLibratoResourceModel {
	if result.Email.IsNull() {
		result.Email = plan.Email
	}
	if result.APIToken.IsNull() {
		result.APIToken = plan.APIToken
	}
	if result.MetricName.IsNull() {
		result.MetricName = plan.MetricName
	}
	return result
}

func (a IntegrationLibratoResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationLibratoResourceModel, diag.Diagnostics) {
	model := *new(IntegrationLibratoResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationLibratoResourceModelAdapter) ToAPIArgument(model IntegrationLibratoResourceModel) (*upapi.IntegrationLibrato, error) {
	return &upapi.IntegrationLibrato{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		Email:         model.Email.ValueString(),
		APIToken:      model.APIToken.ValueString(),
		MetricName:    model.MetricName.ValueString(),
	}, nil
}

func (a IntegrationLibratoResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationLibratoResourceModel, error) {
	return &IntegrationLibratoResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		Email:         types.StringNull(),
		APIToken:      types.StringNull(),
		MetricName:    types.StringNull(),
	}, nil
}

type IntegrationLibratoResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationLibratoResourceAPI) Create(ctx context.Context, arg upapi.IntegrationLibrato) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateLibrato(ctx, arg)
}

func (a IntegrationLibratoResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationLibratoResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationLibrato) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateLibrato(ctx, pk, arg)
}

func (a IntegrationLibratoResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
