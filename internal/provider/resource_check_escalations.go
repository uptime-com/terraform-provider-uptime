package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckEscalationsResource(_ context.Context, p *providerImpl) resource.Resource {
	return &CheckEscalationsResource{
		provider: p,
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_escalations",
			Schema: schema.Schema{
				Description: `Manages escalation rules for a check. Escalations allow you to configure
multiple levels of notifications that are triggered after specific wait times when a check is down.
Each escalation level can send alerts to different contact groups and be repeated multiple times.

Note: This resource manages the escalation configuration for an existing check.
The check must be created first using one of the uptime_check_* resources.`,
				Attributes: map[string]schema.Attribute{
					"check_id": schema.Int64Attribute{
						Required:    true,
						Description: "The ID of the check to configure escalations for",
					},
					"escalations": schema.ListNestedAttribute{
						Required: true,
						Description: `List of escalation rules. Each escalation is triggered sequentially
after the specified wait time. If the list is empty, all escalations will be removed from the check.`,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"wait_time": schema.Int64Attribute{
									Required: true,
									Description: `Time to wait (in seconds) before triggering this escalation level.
For the first escalation, this is the time after the initial alert. For subsequent escalations,
this is the time after the previous escalation.`,
								},
								"num_repeats": schema.Int64Attribute{
									Required: true,
									Description: `Number of times to repeat this escalation level.
Use 0 to repeat indefinitely until the check recovers.`,
								},
								"contact_groups": schema.SetAttribute{
									ElementType: types.StringType,
									Required:    true,
									Description: `List of contact group names to receive notifications for this escalation level.
Each contact group can contain multiple contacts (email addresses, phone numbers, or integrations).`,
								},
							},
						},
					},
				},
			},
		},
		adapter: CheckEscalationsResourceModelAdapter{},
	}
}

type CheckEscalationsResource struct {
	provider *providerImpl
	meta     APIResourceMetadata
	adapter  CheckEscalationsResourceModelAdapter
}

func (r *CheckEscalationsResource) Metadata(_ context.Context, rq resource.MetadataRequest, rs *resource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_" + r.meta.TypeNameSuffix
}

func (r *CheckEscalationsResource) Schema(_ context.Context, _ resource.SchemaRequest, rs *resource.SchemaResponse) {
	rs.Schema = r.meta.Schema
}

func (r *CheckEscalationsResource) Create(ctx context.Context, rq resource.CreateRequest, rs *resource.CreateResponse) {
	model, diags := r.adapter.Get(ctx, rq.Plan)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	arg, err := r.adapter.ToAPIArgument(*model)
	if err != nil {
		rs.Diagnostics.AddError("To API Argument Conversion Failed", err.Error())
		return
	}

	// UpdateEscalations returns *CheckEscalations, which contains Escalations []CheckEscalation
	result, err := r.provider.api.Checks().UpdateEscalations(ctx, *model, upapi.CheckEscalations{Escalations: arg})
	if err != nil {
		rs.Diagnostics.AddError("API Update Escalations Operation Failed", err.Error())
		return
	}

	// result is *CheckEscalations, so result.Escalations is []CheckEscalation
	resultModel, err := r.adapter.FromAPIResult(result.Escalations)
	if err != nil {
		rs.Diagnostics.AddError("From API Result Conversion Failed", err.Error())
		return
	}

	// Preserve the check_id from the plan
	resultModel.CheckID = model.CheckID

	diags = rs.State.Set(ctx, resultModel)
	rs.Diagnostics.Append(diags...)
}

func (r *CheckEscalationsResource) Read(ctx context.Context, rq resource.ReadRequest, rs *resource.ReadResponse) {
	model, diags := r.adapter.Get(ctx, rq.State)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	result, err := r.provider.api.Checks().GetEscalations(ctx, *model)
	if err != nil {
		rs.Diagnostics.AddError("API Get Escalations Operation Failed", err.Error())
		return
	}

	resultModel, err := r.adapter.FromAPIResult(result.Escalations)
	if err != nil {
		rs.Diagnostics.AddError("From API Result Conversion Failed", err.Error())
		return
	}

	// Preserve the check_id from state
	resultModel.CheckID = model.CheckID

	diags = rs.State.Set(ctx, resultModel)
	rs.Diagnostics.Append(diags...)
}

func (r *CheckEscalationsResource) Update(ctx context.Context, rq resource.UpdateRequest, rs *resource.UpdateResponse) {
	model, diags := r.adapter.Get(ctx, rq.Plan)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	arg, err := r.adapter.ToAPIArgument(*model)
	if err != nil {
		rs.Diagnostics.AddError("To API Argument Conversion Failed", err.Error())
		return
	}

	result, err := r.provider.api.Checks().UpdateEscalations(ctx, *model, upapi.CheckEscalations{Escalations: arg})
	if err != nil {
		rs.Diagnostics.AddError("API Update Escalations Operation Failed", err.Error())
		return
	}

	resultModel, err := r.adapter.FromAPIResult(result.Escalations)
	if err != nil {
		rs.Diagnostics.AddError("From API Result Conversion Failed", err.Error())
		return
	}

	// Preserve the check_id from the plan
	resultModel.CheckID = model.CheckID

	diags = rs.State.Set(ctx, resultModel)
	rs.Diagnostics.Append(diags...)
}

func (r *CheckEscalationsResource) Delete(ctx context.Context, rq resource.DeleteRequest, rs *resource.DeleteResponse) {
	model, diags := r.adapter.Get(ctx, rq.State)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	// To delete escalations, we update with an empty list
	_, err := r.provider.api.Checks().UpdateEscalations(ctx, *model, upapi.CheckEscalations{Escalations: []upapi.CheckEscalation{}})
	if err != nil {
		rs.Diagnostics.AddError("API Delete Escalations Operation Failed", err.Error())
		return
	}

	// State is automatically cleared after successful delete
}

type CheckEscalationsResourceModel struct {
	CheckID     types.Int64 `tfsdk:"check_id"`
	Escalations types.List  `tfsdk:"escalations"`

	escalations *escalationsAttribute `tfsdk:"-"`
}

func (m CheckEscalationsResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.CheckID.ValueInt64())
}

type CheckEscalationsResourceModelAdapter struct {
	escalationsAttributeContextAdapter
}

func (a CheckEscalationsResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckEscalationsResourceModel, diag.Diagnostics) {
	model := *new(CheckEscalationsResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.escalations, diags = a.escalationsAttributeContext(ctx, model.Escalations)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckEscalationsResourceModelAdapter) ToAPIArgument(model CheckEscalationsResourceModel) ([]upapi.CheckEscalation, error) {
	if model.escalations == nil {
		// Return empty escalations to remove all escalations
		return []upapi.CheckEscalation{}, nil
	}
	return a.escalationsToAPI(model.escalations), nil
}

func (a CheckEscalationsResourceModelAdapter) FromAPIResult(api []upapi.CheckEscalation) (*CheckEscalationsResourceModel, error) {
	escalations := a.escalationsFromAPI(api)

	var escalationsValue types.List
	if escalations != nil {
		escalationsValue = a.escalationsAttributeValue(*escalations)
	} else {
		// Empty list instead of null
		escalationsValue = a.escalationsAttributeValue(escalationsAttribute{})
	}

	model := CheckEscalationsResourceModel{
		// CheckID needs to be preserved from state - it will be set during CRUD operations
		CheckID:     types.Int64Null(),
		Escalations: escalationsValue,
	}

	return &model, nil
}

// Internal types and methods for escalations handling

type escalationAttribute struct {
	WaitTime      types.Int64 `tfsdk:"wait_time"`
	NumRepeats    types.Int64 `tfsdk:"num_repeats"`
	ContactGroups types.Set   `tfsdk:"contact_groups"`
}

type escalationsAttribute []escalationAttribute

type escalationsAttributeContextAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a escalationsAttributeContextAdapter) escalationAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"wait_time":      types.Int64Type,
		"num_repeats":    types.Int64Type,
		"contact_groups": types.SetType{ElemType: types.StringType},
	}
}

func (a escalationsAttributeContextAdapter) escalationAttributeValues(m escalationAttribute) map[string]attr.Value {
	return map[string]attr.Value{
		"wait_time":      m.WaitTime,
		"num_repeats":    m.NumRepeats,
		"contact_groups": m.ContactGroups,
	}
}

func (a escalationsAttributeContextAdapter) escalationsAttributeContext(ctx context.Context, v types.List) (*escalationsAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	var escalations []escalationAttribute
	diags := v.ElementsAs(ctx, &escalations, false)
	if diags.HasError() {
		return nil, diags
	}

	result := escalationsAttribute(escalations)
	return &result, nil
}

func (a escalationsAttributeContextAdapter) escalationsAttributeValue(escalations escalationsAttribute) types.List {
	if escalations == nil {
		return types.ListNull(types.ObjectType{
			AttrTypes: a.escalationAttributeTypes(),
		})
	}

	elements := make([]attr.Value, len(escalations))
	for i, esc := range escalations {
		elements[i] = types.ObjectValueMust(
			a.escalationAttributeTypes(),
			a.escalationAttributeValues(esc),
		)
	}

	return types.ListValueMust(
		types.ObjectType{AttrTypes: a.escalationAttributeTypes()},
		elements,
	)
}

func (a escalationsAttributeContextAdapter) escalationsToAPI(escalations *escalationsAttribute) []upapi.CheckEscalation {
	if escalations == nil {
		return nil
	}

	result := make([]upapi.CheckEscalation, len(*escalations))
	for i, esc := range *escalations {
		result[i] = upapi.CheckEscalation{
			WaitTime:      int(esc.WaitTime.ValueInt64()),
			NumRepeats:    int(esc.NumRepeats.ValueInt64()),
			ContactGroups: a.ContactGroups(esc.ContactGroups),
		}
	}

	return result
}

func (a escalationsAttributeContextAdapter) escalationsFromAPI(apiEscalations []upapi.CheckEscalation) *escalationsAttribute {
	if len(apiEscalations) == 0 {
		return nil
	}

	result := make(escalationsAttribute, len(apiEscalations))
	for i, esc := range apiEscalations {
		result[i] = escalationAttribute{
			WaitTime:      types.Int64Value(int64(esc.WaitTime)),
			NumRepeats:    types.Int64Value(int64(esc.NumRepeats)),
			ContactGroups: a.ContactGroupsValue(esc.ContactGroups),
		}
	}

	return &result
}
