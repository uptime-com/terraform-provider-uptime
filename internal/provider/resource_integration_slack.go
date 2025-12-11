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

func NewIntegrationSlackResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[IntegrationSlackResourceModel, upapi.IntegrationSlack, upapi.Integration](
		IntegrationSlackResourceAPI{provider: p},
		IntegrationSlackResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "integration_slack",
			Schema: schema.Schema{
				Description: "Slack integration resource. Import using the integration ID: `terraform import uptime_integration_slack.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"webhook_url": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Slack webhook URL",
					},
					"channel": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Slack channel to post to (overrides webhook default)",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type IntegrationSlackResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	WebhookURL    types.String `tfsdk:"webhook_url"`
	Channel       types.String `tfsdk:"channel"`
}

func (m IntegrationSlackResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationSlackResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationSlackResourceModelAdapter) PreservePlanValues(result *IntegrationSlackResourceModel, plan *IntegrationSlackResourceModel) *IntegrationSlackResourceModel {
	if result.WebhookURL.IsNull() {
		result.WebhookURL = plan.WebhookURL
	}
	if result.Channel.IsNull() {
		result.Channel = plan.Channel
	}
	return result
}

func (a IntegrationSlackResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationSlackResourceModel, diag.Diagnostics) {
	model := *new(IntegrationSlackResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationSlackResourceModelAdapter) ToAPIArgument(model IntegrationSlackResourceModel) (*upapi.IntegrationSlack, error) {
	return &upapi.IntegrationSlack{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		WebhookURL:    model.WebhookURL.ValueString(),
		Channel:       model.Channel.ValueString(),
	}, nil
}

func (a IntegrationSlackResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationSlackResourceModel, error) {
	return &IntegrationSlackResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		WebhookURL:    types.StringNull(),
		Channel:       types.StringNull(),
	}, nil
}

type IntegrationSlackResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationSlackResourceAPI) Create(ctx context.Context, arg upapi.IntegrationSlack) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateSlack(ctx, arg)
}

func (a IntegrationSlackResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationSlackResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationSlack) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateSlack(ctx, pk, arg)
}

func (a IntegrationSlackResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
