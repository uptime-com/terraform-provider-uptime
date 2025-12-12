package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationMicrosoftTeamsResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationMicrosoftTeamsResourceModel, upapi.IntegrationMicrosoftTeams, upapi.Integration](
		IntegrationMicrosoftTeamsResourceAPI{provider: p},
		IntegrationMicrosoftTeamsResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_microsoft_teams",
			Schema: schema.Schema{
				Description: "Microsoft Teams integration resource. Import using the integration ID: `terraform import uptime_integration_microsoft_teams.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"webhook_url": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Microsoft Teams webhook URL",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type IntegrationMicrosoftTeamsResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	WebhookURL    types.String `tfsdk:"webhook_url"`
}

func (m IntegrationMicrosoftTeamsResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationMicrosoftTeamsResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationMicrosoftTeamsResourceModelAdapter) PreservePlanValues(result *IntegrationMicrosoftTeamsResourceModel, plan *IntegrationMicrosoftTeamsResourceModel) *IntegrationMicrosoftTeamsResourceModel {
	if result.WebhookURL.IsNull() {
		result.WebhookURL = plan.WebhookURL
	}
	return result
}

func (a IntegrationMicrosoftTeamsResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationMicrosoftTeamsResourceModel, diag.Diagnostics) {
	model := *new(IntegrationMicrosoftTeamsResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationMicrosoftTeamsResourceModelAdapter) ToAPIArgument(model IntegrationMicrosoftTeamsResourceModel) (*upapi.IntegrationMicrosoftTeams, error) {
	return &upapi.IntegrationMicrosoftTeams{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		WebhookUrl:    model.WebhookURL.ValueString(),
	}, nil
}

func (a IntegrationMicrosoftTeamsResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationMicrosoftTeamsResourceModel, error) {
	return &IntegrationMicrosoftTeamsResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		WebhookURL:    types.StringNull(),
	}, nil
}

type IntegrationMicrosoftTeamsResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationMicrosoftTeamsResourceAPI) Create(ctx context.Context, arg upapi.IntegrationMicrosoftTeams) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateMicrosoftTeams(ctx, arg)
}

func (a IntegrationMicrosoftTeamsResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationMicrosoftTeamsResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationMicrosoftTeams) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateMicrosoftTeams(ctx, pk, arg)
}

func (a IntegrationMicrosoftTeamsResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
