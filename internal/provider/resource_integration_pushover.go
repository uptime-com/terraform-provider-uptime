package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationPushoverResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationPushoverResourceModel, upapi.IntegrationPushover, upapi.Integration](
		IntegrationPushoverResourceAPI{provider: p},
		IntegrationPushoverResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_pushover",
			Schema: schema.Schema{
				Description: "Pushover integration resource. Import using the integration ID: `terraform import uptime_integration_pushover.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"user": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Pushover user key",
					},
					"priority": schema.Int64Attribute{
						Required:    true,
						Description: "Message priority (-2 to 2)",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type IntegrationPushoverResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	User          types.String `tfsdk:"user"`
	Priority      types.Int64  `tfsdk:"priority"`
}

func (m IntegrationPushoverResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationPushoverResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationPushoverResourceModelAdapter) PreservePlanValues(result *IntegrationPushoverResourceModel, plan *IntegrationPushoverResourceModel) *IntegrationPushoverResourceModel {
	if result.User.IsNull() {
		result.User = plan.User
	}
	if result.Priority.IsNull() {
		result.Priority = plan.Priority
	}
	return result
}

func (a IntegrationPushoverResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationPushoverResourceModel, diag.Diagnostics) {
	model := *new(IntegrationPushoverResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationPushoverResourceModelAdapter) ToAPIArgument(model IntegrationPushoverResourceModel) (*upapi.IntegrationPushover, error) {
	return &upapi.IntegrationPushover{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		User:          model.User.ValueString(),
		Priority:      model.Priority.ValueInt64(),
	}, nil
}

func (a IntegrationPushoverResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationPushoverResourceModel, error) {
	return &IntegrationPushoverResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		User:          types.StringNull(),
		Priority:      types.Int64Null(),
	}, nil
}

type IntegrationPushoverResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationPushoverResourceAPI) Create(ctx context.Context, arg upapi.IntegrationPushover) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreatePushover(ctx, arg)
}

func (a IntegrationPushoverResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationPushoverResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationPushover) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdatePushover(ctx, pk, arg)
}

func (a IntegrationPushoverResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
