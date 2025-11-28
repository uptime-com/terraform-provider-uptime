package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewIntegrationOpsgenieResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[IntegrationOpsgenieResourceModel, upapi.IntegrationOpsgenie, upapi.Integration]{
		api: IntegrationOpsgenieResourceAPI{provider: p},
		mod: IntegrationOpsgenieResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "integration_opsgenie",
			Schema: schema.Schema{
				Description: "Opsgenie integration resource",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"api_endpoint": schema.StringAttribute{
						Required: true,
					},
					"api_key": schema.StringAttribute{
						Required: true,
					},
					"teams": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "A comma separated list of team names which will be responsible for the alert",
					},
					"tags": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
						Description: ("A comma separated list of labels attached to the alert. " +
							"You may overwrite the quiet hours setting for urgent alerts by adding the OverwriteQuietHours tag. " +
							"Leave blank to automatically pull the tags from the check instead."),
					},
					"auto_resolve": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Automatically resolve incident once the check is back up.",
						Default:     booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

type IntegrationOpsgenieResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	APIEndpoint   types.String `tfsdk:"api_endpoint"`
	APIKey        types.String `tfsdk:"api_key"`
	Teams         types.String `tfsdk:"teams"`
	Tags          types.String `tfsdk:"tags"`
	AutoResolve   types.Bool   `tfsdk:"auto_resolve"`
}

func (m IntegrationOpsgenieResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type IntegrationOpsgenieResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a IntegrationOpsgenieResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*IntegrationOpsgenieResourceModel, diag.Diagnostics) {
	model := *new(IntegrationOpsgenieResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a IntegrationOpsgenieResourceModelAdapter) ToAPIArgument(model IntegrationOpsgenieResourceModel) (*upapi.IntegrationOpsgenie, error) {
	return &upapi.IntegrationOpsgenie{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroupsSlice(model.ContactGroups),
		APIEndpoint:   model.APIEndpoint.ValueString(),
		APIKey:        model.APIKey.ValueString(),
		Teams:         model.Teams.ValueString(),
		Tags:          model.Tags.ValueString(),
		Autoresolve:   model.AutoResolve.ValueBool(),
	}, nil
}

func (a IntegrationOpsgenieResourceModelAdapter) FromAPIResult(api upapi.Integration) (*IntegrationOpsgenieResourceModel, error) {
	return &IntegrationOpsgenieResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsSliceValue(api.ContactGroups),
		APIEndpoint:   types.StringValue(api.APIEndpoint),
		APIKey:        types.StringValue(api.APIKey),
		Teams:         types.StringValue(api.Teams),
		Tags:          types.StringValue(api.Tags),
		AutoResolve:   types.BoolValue(api.Autoresolve),
	}, nil
}

type IntegrationOpsgenieResourceAPI struct {
	provider *providerImpl
}

func (a IntegrationOpsgenieResourceAPI) Create(ctx context.Context, arg upapi.IntegrationOpsgenie) (*upapi.Integration, error) {
	return a.provider.api.Integrations().CreateOpsgenie(ctx, arg)
}

func (a IntegrationOpsgenieResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Integration, error) {
	return a.provider.api.Integrations().Get(ctx, pk)
}

func (a IntegrationOpsgenieResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.IntegrationOpsgenie) (*upapi.Integration, error) {
	return a.provider.api.Integrations().UpdateOpsgenie(ctx, pk, arg)
}

func (a IntegrationOpsgenieResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Integrations().Delete(ctx, pk)
}
