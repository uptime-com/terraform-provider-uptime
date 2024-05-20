package provider

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckRUM2Resource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckRUM2ResourceModel, upapi.CheckRUM2, upapi.Check]{
		api: CheckRUM2ResourceAPI{provider: p},
		mod: CheckRUM2ResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_rum2",
			Schema: schema.Schema{
				Description: "Create a new Real User Monitoring check",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"address":                   AddressHostnameSchemaAttribute(),
					"sla_uptime":                SLAUptimeSchemaAttribute(),
					"notes":                     NotesSchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
				},
			},
		},
	}
}

type CheckRUM2ResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	URL                    types.String `tfsdk:"url"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Address                types.String `tfsdk:"address"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLAUptime              Decimal      `tfsdk:"sla_uptime"`
}

func (m CheckRUM2ResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckRUM2ResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	TagsAttributeAdapter
}

func (a CheckRUM2ResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckRUM2ResourceModel, diag.Diagnostics) {
	model := *new(CheckRUM2ResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckRUM2ResourceModelAdapter) ToAPIArgument(model CheckRUM2ResourceModel) (*upapi.CheckRUM2, error) {
	api := upapi.CheckRUM2{
		Name:                   model.Name.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Address:                model.Address.ValueString(),
		Notes:                  model.Notes.ValueString(),
		UptimeSLA:              model.SLAUptime.ValueDecimal(),
		IncludeInGlobalMetrics: model.IncludeInGlobalMetrics.ValueBool(),
	}

	return &api, nil
}

func (a CheckRUM2ResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckRUM2ResourceModel, error) {
	merr := new(multierror.Error)
	model := CheckRUM2ResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Address:                types.StringValue(api.Address),
		Notes:                  types.StringValue(api.Notes),
		SLAUptime:              DecimalValue(api.UptimeSLA),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
	}
	return &model, merr.ErrorOrNil()
}

type CheckRUM2ResourceAPI struct {
	provider *providerImpl
}

func (a CheckRUM2ResourceAPI) Create(ctx context.Context, arg upapi.CheckRUM2) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateRUM2(ctx, arg)
}

func (a CheckRUM2ResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a CheckRUM2ResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckRUM2) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateRUM2(ctx, pk, arg)
}

func (a CheckRUM2ResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
