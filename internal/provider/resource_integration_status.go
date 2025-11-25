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

func NewIntegrationStatusResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationStatusResourceModel, upapi.IntegrationStatus, upapi.Integration]{
		api: IntegrationStatusResourceAPI{provider: p},
		mod: IntegrationStatusResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_status",
			Schema: schema.Schema{
				Description: "Status.io integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"statuspage_id": schema.StringAttribute{
						Required:    true,
						Description: "Status.io status page ID",
					},
					"api_id": schema.StringAttribute{
						Required:    true,
						Description: "Status.io API ID",
					},
					"api_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Status.io API key",
					},
					"component": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Component ID to update",
					},
					"container": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "Container ID",
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

type IntegrationStatusResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	StatuspageID  types.String `tfsdk:"statuspage_id"`
	APIID         types.String `tfsdk:"api_id"`
	APIKey        types.String `tfsdk:"api_key"`
	Component     types.String `tfsdk:"component"`
	Container     types.String `tfsdk:"container"`
	Metric        types.String `tfsdk:"metric"`
}

func (m IntegrationStatusResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationStatusResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationStatusResourceModelAdapter) PreservePlanValues(result *IntegrationStatusResourceModel, plan *IntegrationStatusResourceModel) *IntegrationStatusResourceModel {
	if result.StatuspageID.IsNull() {
		result.StatuspageID = plan.StatuspageID
	}
	if result.APIID.IsNull() {
		result.APIID = plan.APIID
	}
	if result.APIKey.IsNull() {
		result.APIKey = plan.APIKey
	}
	if result.Component.IsNull() {
		result.Component = plan.Component
	}
	if result.Container.IsNull() {
		result.Container = plan.Container
	}
	if result.Metric.IsNull() {
		result.Metric = plan.Metric
	}
	return result
}

func (a IntegrationStatusResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationStatusResourceModel, diag.Diagnostics) {
	model := *new(IntegrationStatusResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationStatusResourceModelAdapter) ToAPIArgument(model IntegrationStatusResourceModel) (*upapi.IntegrationStatus, error) {
	return &upapi.IntegrationStatus{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroups(model.ContactGroups),
		StatuspageID:  model.StatuspageID.ValueString(),
		APIID:         model.APIID.ValueString(),
		APIKey:        model.APIKey.ValueString(),
		Component:     model.Component.ValueString(),
		Container:     model.Container.ValueString(),
		Metric:        model.Metric.ValueString(),
	}, nil
}

func (a IntegrationStatusResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationStatusResourceModel, error) {
	return &IntegrationStatusResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsValue(api.ContactGroups),
		StatuspageID:  types.StringNull(),
		APIID:         types.StringNull(),
		APIKey:        types.StringNull(),
		Component:     types.StringNull(),
		Container:     types.StringNull(),
		Metric:        types.StringNull(),
	}, nil
}

type IntegrationStatusResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationStatusResourceAPI) Create(ctx context.Context, arg upapi.IntegrationStatus) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateStatus(ctx, arg)
}

func (a IntegrationStatusResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationStatusResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationStatus) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateStatus(ctx, pk, arg)
}

func (a IntegrationStatusResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
