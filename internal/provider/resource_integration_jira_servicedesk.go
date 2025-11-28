package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationJiraServicedeskResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationJiraServicedeskResourceModel, upapi.IntegrationJiraServicedesk, upapi.Integration]{
		api: IntegrationJiraServicedeskResourceAPI{provider: p},
		mod: IntegrationJiraServicedeskResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_jira_servicedesk",
			Schema: schema.Schema{
				Description: "JIRA Service Desk integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"api_email": schema.StringAttribute{
						Required:    true,
						Description: "Email address for JIRA API authentication",
					},
					"api_token": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "API token for JIRA authentication",
					},
					"jira_subdomain": schema.StringAttribute{
						Required:    true,
						Description: "JIRA subdomain (e.g., 'mycompany' for mycompany.atlassian.net)",
					},
					"project_key": schema.StringAttribute{
						Required:    true,
						Description: "JIRA project key",
					},
					"labels": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Comma-separated list of labels to add to created issues",
					},
					"custom_field_id_account_name": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
						Description: "Custom field ID for account name",
					},
					"custom_field_id_check_name": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
						Description: "Custom field ID for check name",
					},
					"custom_field_id_check_url": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
						Description: "Custom field ID for check URL",
					},
					"custom_fields_json": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Additional custom fields as JSON",
					},
				},
			},
		},
	}
}

type IntegrationJiraServicedeskResourceModel struct {
	ID                       types.Int64  `tfsdk:"id"`
	URL                      types.String `tfsdk:"url"`
	Name                     types.String `tfsdk:"name"`
	ContactGroups            types.Set    `tfsdk:"contact_groups"`
	APIEmail                 types.String `tfsdk:"api_email"`
	APIToken                 types.String `tfsdk:"api_token"`
	JiraSubdomain            types.String `tfsdk:"jira_subdomain"`
	ProjectKey               types.String `tfsdk:"project_key"`
	Labels                   types.String `tfsdk:"labels"`
	CustomFieldIdAccountName types.Int64  `tfsdk:"custom_field_id_account_name"`
	CustomFieldIdCheckName   types.Int64  `tfsdk:"custom_field_id_check_name"`
	CustomFieldIdCheckUrl    types.Int64  `tfsdk:"custom_field_id_check_url"`
	CustomFieldsJson         types.String `tfsdk:"custom_fields_json"`
}

func (m IntegrationJiraServicedeskResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationJiraServicedeskResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationJiraServicedeskResourceModelAdapter) PreservePlanValues(result *IntegrationJiraServicedeskResourceModel, plan *IntegrationJiraServicedeskResourceModel) *IntegrationJiraServicedeskResourceModel {
	if result.APIEmail.IsNull() {
		result.APIEmail = plan.APIEmail
	}
	if result.APIToken.IsNull() {
		result.APIToken = plan.APIToken
	}
	if result.JiraSubdomain.IsNull() {
		result.JiraSubdomain = plan.JiraSubdomain
	}
	if result.ProjectKey.IsNull() {
		result.ProjectKey = plan.ProjectKey
	}
	if result.Labels.IsNull() {
		result.Labels = plan.Labels
	}
	if result.CustomFieldIdAccountName.IsNull() {
		result.CustomFieldIdAccountName = plan.CustomFieldIdAccountName
	}
	if result.CustomFieldIdCheckName.IsNull() {
		result.CustomFieldIdCheckName = plan.CustomFieldIdCheckName
	}
	if result.CustomFieldIdCheckUrl.IsNull() {
		result.CustomFieldIdCheckUrl = plan.CustomFieldIdCheckUrl
	}
	if result.CustomFieldsJson.IsNull() {
		result.CustomFieldsJson = plan.CustomFieldsJson
	}
	return result
}

func (a IntegrationJiraServicedeskResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationJiraServicedeskResourceModel, diag.Diagnostics) {
	model := *new(IntegrationJiraServicedeskResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationJiraServicedeskResourceModelAdapter) ToAPIArgument(model IntegrationJiraServicedeskResourceModel) (*upapi.IntegrationJiraServicedesk, error) {
	return &upapi.IntegrationJiraServicedesk{
		Name:                     model.Name.ValueString(),
		ContactGroups:            a.ContactGroupsSlice(model.ContactGroups),
		APIEmail:                 model.APIEmail.ValueString(),
		APIToken:                 model.APIToken.ValueString(),
		JiraSubdomain:            model.JiraSubdomain.ValueString(),
		ProjectKey:               model.ProjectKey.ValueString(),
		Labels:                   model.Labels.ValueString(),
		CustomFieldIdAccountName: model.CustomFieldIdAccountName.ValueInt64(),
		CustomFieldIdCheckName:   model.CustomFieldIdCheckName.ValueInt64(),
		CustomFieldIdCheckUrl:    model.CustomFieldIdCheckUrl.ValueInt64(),
		CustomFieldsJson:         model.CustomFieldsJson.ValueString(),
	}, nil
}

func (a IntegrationJiraServicedeskResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationJiraServicedeskResourceModel, error) {
	return &IntegrationJiraServicedeskResourceModel{
		ID:                       types.Int64Value(api.PK),
		URL:                      types.StringValue(api.URL),
		Name:                     types.StringValue(api.Name),
		ContactGroups:            a.ContactGroupsSliceValue(api.ContactGroups),
		APIEmail:                 types.StringNull(),
		APIToken:                 types.StringNull(),
		JiraSubdomain:            types.StringNull(),
		ProjectKey:               types.StringNull(),
		Labels:                   types.StringNull(),
		CustomFieldIdAccountName: types.Int64Null(),
		CustomFieldIdCheckName:   types.Int64Null(),
		CustomFieldIdCheckUrl:    types.Int64Null(),
		CustomFieldsJson:         types.StringNull(),
	}, nil
}

type IntegrationJiraServicedeskResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationJiraServicedeskResourceAPI) Create(ctx context.Context, arg upapi.IntegrationJiraServicedesk) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateJiraServicedesk(ctx, arg)
}

func (a IntegrationJiraServicedeskResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationJiraServicedeskResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationJiraServicedesk) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateJiraServiceDesk(ctx, pk, arg)
}

func (a IntegrationJiraServicedeskResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
