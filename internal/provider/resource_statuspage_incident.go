package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageIncidentResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[StatusPageIncidentResourceModel, StatusPageIncidentWrapper, StatusPageIncidentWrapper]{
		api: &StatusPageIncidentResourceAPI{provider: p},
		mod: StatusPageIncidentResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "statuspage_incident",
			Schema: schema.Schema{
				Description: "Status page incident or maintenance window resource",
				Attributes: map[string]schema.Attribute{
					"statuspage_id": schema.Int64Attribute{
						Required: true,
					},
					"id":   IDSchemaAttribute(),
					"url":  URLSchemaAttribute(),
					"name": NameSchemaAttribute(),
					"starts_at": schema.StringAttribute{
						Description: "When this incident occurred in GMT",
						Required:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"ends_at": schema.StringAttribute{
						Description: "When this incident ended in GMT",
						Optional:    true,
						Computed:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"include_in_global_metrics": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"updates": schema.SetNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": IDSchemaAttribute(),
								"description": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"incident_state": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										OneOfStringValidator([]string{"investigating", "identified", "monitoring", "resolved", "notification", "maintenance"}),
									},
								},
							},
						},
					},
					"affected_components": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": IDSchemaAttribute(),
								"status": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										OneOfStringValidator([]string{"major-outage", "partial-outage", "degraded-performance", "under-maintenance"}),
									},
								},
								"component_id": schema.Int64Attribute{
									Required: true,
								},
							},
						},
					},
					"incident_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("major-outage"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"INCIDENT", "SCHEDULED_MAINTENANCE"}),
						},
					},
					"update_component_status": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"notify_subscribers": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"send_maintenance_start_notification": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

type StatusPageIncidentWrapper struct {
	upapi.StatusPageIncident
	StatusPageID int64
}

func (w StatusPageIncidentWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.PK)
}

type StatusPageIncidentResourceModel struct {
	StatusPageID                     types.Int64       `tfsdk:"statuspage_id"`
	ID                               types.Int64       `tfsdk:"id"`
	URL                              types.String      `tfsdk:"url"`
	Name                             types.String      `tfsdk:"name"`
	EndsAt                           timetypes.RFC3339 `tfsdk:"ends_at"`
	StartsAt                         timetypes.RFC3339 `tfsdk:"starts_at"`
	IncludeInGlobalMetrics           types.Bool        `tfsdk:"include_in_global_metrics"`
	Updates                          types.Set         `tfsdk:"updates"`
	AffectedComponents               types.Set         `tfsdk:"affected_components"`
	IncidentType                     types.String      `tfsdk:"incident_type"`
	UpdateComponentStatus            types.Bool        `tfsdk:"update_component_status"`
	NotifySubscribers                types.Bool        `tfsdk:"notify_subscribers"`
	SendMaintenanceStartNotification types.Bool        `tfsdk:"send_maintenance_start_notification"`

	updates            []StatusPageIncidentUpdateAttribute
	affectedComponents []StatusPageIncidentAffectedComponentAttribute
}

func (m StatusPageIncidentResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageIncidentResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a StatusPageIncidentResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*StatusPageIncidentResourceModel, diag.Diagnostics) {
	var model StatusPageIncidentResourceModel
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}

	model.updates, diags = a.UpdatesContext(ctx, model.Updates)
	if diags.HasError() {
		return nil, diags
	}

	model.affectedComponents, diags = a.AffectedComponentsContext(ctx, model.AffectedComponents)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a StatusPageIncidentResourceModelAdapter) ToAPIArgument(
	model StatusPageIncidentResourceModel,
) (*StatusPageIncidentWrapper, error) {
	api := StatusPageIncidentWrapper{
		StatusPageID: model.StatusPageID.ValueInt64(),
		StatusPageIncident: upapi.StatusPageIncident{
			Name:                             model.Name.ValueString(),
			IncludeInGlobalMetrics:           model.IncludeInGlobalMetrics.ValueBool(),
			IncidentType:                     model.IncidentType.ValueString(),
			UpdateComponentStatus:            model.UpdateComponentStatus.ValueBool(),
			NotifySubscribers:                model.NotifySubscribers.ValueBool(),
			SendMaintenanceStartNotification: model.SendMaintenanceStartNotification.ValueBool(),
		},
	}

	if !model.EndsAt.IsNull() && !model.EndsAt.IsUnknown() {
		api.EndsAt = model.EndsAt.ValueString()
	}

	if !model.StartsAt.IsNull() && !model.StartsAt.IsUnknown() {
		api.StartsAt = model.StartsAt.ValueString()
	}

	if len(model.updates) != 0 {
		updates := make([]upapi.IncidentUpdate, 0)
		for _, v := range model.updates {
			updates = append(updates, upapi.IncidentUpdate{
				Description:   v.Description.ValueString(),
				IncidentState: v.IncidentState.ValueString(),
			})
		}
		api.Updates = updates
	}

	if len(model.affectedComponents) != 0 {
		affectedComponents := make([]upapi.IncidentAffectedComponentEntity, 0)
		for _, v := range model.affectedComponents {
			affectedComponents = append(affectedComponents, upapi.IncidentAffectedComponentEntity{
				PK:     v.ID.ValueInt64(),
				Status: v.Status.ValueString(),
				Component: upapi.IncidentAffectedComponent{
					PK: v.ComponentID.ValueInt64(),
				},
			})
		}
		api.AffectedComponents = affectedComponents
	}

	return &api, nil
}

func (a StatusPageIncidentResourceModelAdapter) FromAPIResult(
	api StatusPageIncidentWrapper,
) (*StatusPageIncidentResourceModel, error) {
	model := &StatusPageIncidentResourceModel{
		StatusPageID:                     types.Int64Value(api.StatusPageID),
		ID:                               types.Int64Value(api.PK),
		URL:                              types.StringValue(api.URL),
		Name:                             types.StringValue(api.Name),
		IncludeInGlobalMetrics:           types.BoolValue(api.IncludeInGlobalMetrics),
		IncidentType:                     types.StringValue(api.IncidentType),
		UpdateComponentStatus:            types.BoolValue(api.UpdateComponentStatus),
		NotifySubscribers:                types.BoolValue(api.NotifySubscribers),
		SendMaintenanceStartNotification: types.BoolValue(api.SendMaintenanceStartNotification),
	}

	var d diag.Diagnostics
	if api.EndsAt != "" {
		model.EndsAt, d = timetypes.NewRFC3339PointerValue(&api.EndsAt)
		if d.HasError() {
			return nil, fmt.Errorf("error parsing EndsAt: %v", d)
		}
	}

	if api.StartsAt != "" {
		model.StartsAt, d = timetypes.NewRFC3339PointerValue(&api.StartsAt)
		if d.HasError() {
			return nil, fmt.Errorf("error parsing StartsAt: %v", d)
		}
	}

	updates := []StatusPageIncidentUpdateAttribute{}
	for _, item := range api.Updates {
		updates = append(updates, StatusPageIncidentUpdateAttribute{
			ID:            types.Int64Value(item.PK),
			Description:   types.StringValue(item.Description),
			IncidentState: types.StringValue(item.IncidentState),
		})
	}

	var diags diag.Diagnostics
	if model.Updates, diags = a.UpdatesValue(updates); diags.HasError() {
		return nil, fmt.Errorf("failed to convert updates: %v", diags)
	}

	affectedComponents := []StatusPageIncidentAffectedComponentAttribute{}
	for _, item := range api.AffectedComponents {
		affectedComponents = append(affectedComponents, StatusPageIncidentAffectedComponentAttribute{
			ID:          types.Int64Value(item.PK),
			Status:      types.StringValue(item.Status),
			ComponentID: types.Int64Value(item.Component.PK),
		})
	}

	if model.AffectedComponents, diags = a.AffectedComponentsValue(affectedComponents); diags.HasError() {
		return nil, fmt.Errorf("failed to convert affected components: %v", diags)
	}

	return model, nil
}

func (a StatusPageIncidentResourceModelAdapter) UpdatesContext(
	ctx context.Context, v types.Set,
) ([]StatusPageIncidentUpdateAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	out := make([]StatusPageIncidentUpdateAttribute, 0)
	if d := v.ElementsAs(ctx, &out, false); d.HasError() {
		return nil, d
	}
	return out, nil
}

func (a StatusPageIncidentResourceModelAdapter) UpdatesValue(
	model []StatusPageIncidentUpdateAttribute,
) (types.Set, diag.Diagnostics) {
	values, diags := a.updatesAttributeValues(model)
	if diags.HasError() {
		return types.Set{}, diags
	}
	return types.SetValueMust(
		types.ObjectType{}.WithAttributeTypes(a.updatesAttributeTypes()), values), diags
}

func (a StatusPageIncidentResourceModelAdapter) updatesAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.Int64Type,
		"description":    types.StringType,
		"incident_state": types.StringType,
	}
}

func (a StatusPageIncidentResourceModelAdapter) updatesAttributeValues(
	model []StatusPageIncidentUpdateAttribute,
) (out []attr.Value, diags diag.Diagnostics) {
	out = make([]attr.Value, len(model))
	for i := range model {
		out[i], diags = types.ObjectValue(a.updatesAttributeTypes(), map[string]attr.Value{
			"id":             model[i].ID,
			"description":    model[i].Description,
			"incident_state": model[i].IncidentState,
		})
		if diags.HasError() {
			return
		}
	}
	return
}

func (a StatusPageIncidentResourceModelAdapter) AffectedComponentsContext(
	ctx context.Context, v types.Set,
) ([]StatusPageIncidentAffectedComponentAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	out := make([]StatusPageIncidentAffectedComponentAttribute, 0)
	if d := v.ElementsAs(ctx, &out, false); d.HasError() {
		return nil, d
	}
	return out, nil
}

func (a StatusPageIncidentResourceModelAdapter) AffectedComponentsValue(
	model []StatusPageIncidentAffectedComponentAttribute,
) (types.Set, diag.Diagnostics) {
	values, diags := a.affectedComponentsAttributeValues(model)
	if diags.HasError() {
		return types.Set{}, diags
	}
	return types.SetValueMust(
		types.ObjectType{}.WithAttributeTypes(a.affectedComponentsAttributeTypes()), values), diags
}

func (a StatusPageIncidentResourceModelAdapter) affectedComponentsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.Int64Type,
		"status":       types.StringType,
		"component_id": types.Int64Type,
	}
}

func (a StatusPageIncidentResourceModelAdapter) affectedComponentsAttributeValues(
	model []StatusPageIncidentAffectedComponentAttribute,
) (out []attr.Value, diags diag.Diagnostics) {
	out = make([]attr.Value, len(model))
	for i := range model {
		out[i], diags = types.ObjectValue(a.affectedComponentsAttributeTypes(), map[string]attr.Value{
			"id":           model[i].ID,
			"status":       model[i].Status,
			"component_id": model[i].ComponentID,
		})
		if diags.HasError() {
			return
		}
	}
	return
}

type StatusPageIncidentResourceAPI struct {
	provider *providerImpl
}

func (a StatusPageIncidentResourceAPI) Create(ctx context.Context, arg StatusPageIncidentWrapper) (*StatusPageIncidentWrapper, error) {
	obj, err := a.provider.api.StatusPages().Incidents(upapi.PrimaryKey(arg.StatusPageID)).Create(ctx, arg.StatusPageIncident)
	if err != nil {
		return nil, err
	}

	return &StatusPageIncidentWrapper{StatusPageIncident: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageIncidentResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*StatusPageIncidentWrapper, error) {
	model, ok := arg.(StatusPageIncidentResourceModel)
	if !ok {
		return nil, fmt.Errorf("resource read failed:unexpected type %T", arg)
	}

	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	obj, err := a.provider.api.StatusPages().Incidents(statusPageID).Get(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &StatusPageIncidentWrapper{StatusPageIncident: *obj, StatusPageID: int64(statusPageID)}, nil
}

func (a StatusPageIncidentResourceAPI) Update(
	ctx context.Context, pk upapi.PrimaryKeyable, arg StatusPageIncidentWrapper,
) (*StatusPageIncidentWrapper, error) {
	obj, err := a.provider.api.StatusPages().
		Incidents(upapi.PrimaryKey(arg.StatusPageID)).
		Update(ctx, pk, arg.StatusPageIncident)
	if err != nil {
		return nil, err
	}
	return &StatusPageIncidentWrapper{StatusPageIncident: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageIncidentResourceAPI) Delete(ctx context.Context, arg upapi.PrimaryKeyable) error {
	model, ok := arg.(StatusPageIncidentResourceModel)
	if !ok {
		return fmt.Errorf("resource delete failed: unexpected type %T", arg)
	}
	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	return a.provider.api.StatusPages().Incidents(statusPageID).Delete(ctx, arg)
}

type StatusPageIncidentUpdateAttribute struct {
	ID            types.Int64  `tfsdk:"id"`
	Description   types.String `tfsdk:"description"`
	IncidentState types.String `tfsdk:"incident_state"`
}

type StatusPageIncidentAffectedComponentAttribute struct {
	ID          types.Int64  `tfsdk:"id"`
	Status      types.String `tfsdk:"status"`
	ComponentID types.Int64  `tfsdk:"component_id"`
}
