package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckMaintenanceResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckMaintenanceResourceModel, CheckMaintenanceWrapper, CheckMaintenanceWrapper]{
		api: &CheckMaintenanceResourceAPI{provider: p},
		mod: CheckMaintenanceResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_maintenance",
			Schema: schema.Schema{
				Description: "Set maintenance windows for a check",
				Attributes: map[string]schema.Attribute{
					"check_id": schema.Int64Attribute{
						Required: true,
					},
					"state": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("ACTIVE"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"SUPPRESSED", "ACTIVE", "SCHEDULED"}),
						},
					},
					"schedule": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										OneOfStringValidator([]string{"WEEKLY", "MONTHLY", "ONCE"}),
									},
								},
								"from_time": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`),
											"Time must be in HH:MM format, 00:00:00 - 23:59:59",
										),
									},
								},
								"to_time": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`),
											"Time must be in HH:MM format, 00:00:00 - 23:59:59",
										),
									},
								},
								"weekdays": schema.SetAttribute{
									Optional:    true,
									Computed:    true,
									ElementType: types.Int32Type,
									Validators: []validator.Set{
										setvalidator.SizeBetween(0, 7),
										setvalidator.ValueInt32sAre(
											int32validator.Between(0, 6),
										),
									},
								},
								"monthday": schema.Int32Attribute{
									Optional: true,
									Computed: true,
									Default:  int32default.StaticInt32(0),
									Validators: []validator.Int32{
										int32validator.Between(0, 30),
									},
								},
								"monthday_from": schema.Int32Attribute{
									Optional: true,
									Computed: true,
									Default:  int32default.StaticInt32(0),
									Validators: []validator.Int32{
										int32validator.Between(0, 30),
									},
								},
								"monthday_to": schema.Int32Attribute{
									Optional: true,
									Computed: true,
									Default:  int32default.StaticInt32(0),
									Validators: []validator.Int32{
										int32validator.Between(0, 30),
									},
								},
								"once_start_date": schema.StringAttribute{
									Optional:   true,
									CustomType: timetypes.RFC3339Type{},
								},
								"once_end_date": schema.StringAttribute{
									Optional:   true,
									CustomType: timetypes.RFC3339Type{},
								},
							},
						},
					},
				},
			},
		},
	}
}

type CheckMaintenanceWrapper struct {
	upapi.CheckMaintenance

	CheckID int64
}

func (w CheckMaintenanceWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.CheckID)
}

type CheckMaintenanceResourceModel struct {
	CheckID  types.Int64  `tfsdk:"check_id"`
	State    types.String `tfsdk:"state"`
	Schedule types.Set    `tfsdk:"schedule"`

	schedule []CheckMaintenanceScheduleAttribute `tfsdk:"-"`
}

func (m CheckMaintenanceResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.CheckID.ValueInt64())
}

type CheckMaintenanceScheduleAttribute struct {
	Type          types.String      `tfsdk:"type"`
	FromTime      types.String      `tfsdk:"from_time"`
	ToTime        types.String      `tfsdk:"to_time"`
	Weekdays      types.Set         `tfsdk:"weekdays"`
	Monthday      types.Int32       `tfsdk:"monthday"`
	MonthdayFrom  types.Int32       `tfsdk:"monthday_from"`
	MonthdayTo    types.Int32       `tfsdk:"monthday_to"`
	OnceStartDate timetypes.RFC3339 `tfsdk:"once_start_date"`
	OnceEndDate   timetypes.RFC3339 `tfsdk:"once_end_date"`
}

type CheckMaintenanceResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a CheckMaintenanceResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*CheckMaintenanceResourceModel, diag.Diagnostics) {
	var model CheckMaintenanceResourceModel
	if diags := sg.Get(ctx, &model); diags.HasError() {
		return nil, diags
	}

	var diags diag.Diagnostics
	model.schedule, diags = a.ScheduleContext(ctx, model.Schedule)
	if diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a CheckMaintenanceResourceModelAdapter) ScheduleContext(ctx context.Context, v types.Set) ([]CheckMaintenanceScheduleAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	out := make([]CheckMaintenanceScheduleAttribute, 0)
	if d := v.ElementsAs(ctx, &out, false); d.HasError() {
		return nil, d
	}
	return out, nil
}

func (a CheckMaintenanceResourceModelAdapter) ScheduleValue(model []CheckMaintenanceScheduleAttribute) (types.Set, diag.Diagnostics) {
	values, diags := a.scheduleAttributeValues(model)
	if diags.HasError() {
		return types.Set{}, diags
	}
	return types.SetValueMust(types.ObjectType{}.WithAttributeTypes(a.scheduleAttributeTypes()), values), diags
}

func (a CheckMaintenanceResourceModelAdapter) scheduleAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":            types.StringType,
		"from_time":       types.StringType,
		"to_time":         types.StringType,
		"weekdays":        types.SetType{}.WithElementType(types.Int32Type),
		"monthday":        types.Int32Type,
		"monthday_from":   types.Int32Type,
		"monthday_to":     types.Int32Type,
		"once_start_date": timetypes.RFC3339Type{},
		"once_end_date":   timetypes.RFC3339Type{},
	}
}

func (a CheckMaintenanceResourceModelAdapter) scheduleAttributeValues(model []CheckMaintenanceScheduleAttribute) (out []attr.Value, diags diag.Diagnostics) {
	out = make([]attr.Value, len(model))
	for i := range model {
		out[i], diags = types.ObjectValue(a.scheduleAttributeTypes(), map[string]attr.Value{
			"type":            model[i].Type,
			"from_time":       model[i].FromTime,
			"to_time":         model[i].ToTime,
			"weekdays":        model[i].Weekdays,
			"monthday":        model[i].Monthday,
			"monthday_from":   model[i].MonthdayFrom,
			"monthday_to":     model[i].MonthdayTo,
			"once_start_date": model[i].OnceStartDate,
			"once_end_date":   model[i].OnceEndDate,
		})
		if diags.HasError() {
			return
		}
	}
	return
}

func (a CheckMaintenanceResourceModelAdapter) ToAPIArgument(
	model CheckMaintenanceResourceModel,
) (*CheckMaintenanceWrapper, error) {
	api := upapi.CheckMaintenance{
		State: model.State.ValueString(),
	}
	if len(model.schedule) > 0 {
		for _, v := range model.schedule {
			s := upapi.CheckMaintenanceSchedule{
				Type:         v.Type.ValueString(),
				FromTime:     v.FromTime.ValueString(),
				ToTime:       v.ToTime.ValueString(),
				Monthday:     int(v.Monthday.ValueInt32()),
				MonthdayFrom: int(v.MonthdayFrom.ValueInt32()),
				MonthdayTo:   int(v.MonthdayTo.ValueInt32()),
			}

			if !v.OnceStartDate.IsNull() && !v.OnceStartDate.IsUnknown() {
				s.OnceStartDate = v.OnceStartDate.ValueString()
			}
			if !v.OnceEndDate.IsNull() && !v.OnceEndDate.IsUnknown() {
				s.OnceEndDate = v.OnceEndDate.ValueString()
			}
			for _, v := range a.Slice(v.Weekdays) {
				s.Weekdays = append(s.Weekdays, int(v))
			}
			api.Schedule = append(api.Schedule, s)
		}
	}

	return &CheckMaintenanceWrapper{
		CheckMaintenance: api,
		CheckID:          model.CheckID.ValueInt64(),
	}, nil
}

func (a CheckMaintenanceResourceModelAdapter) FromAPIResult(
	api CheckMaintenanceWrapper,
) (*CheckMaintenanceResourceModel, error) {
	model := CheckMaintenanceResourceModel{
		CheckID: types.Int64Value(api.CheckID),
		State:   types.StringValue(api.State),
	}

	schedule := []CheckMaintenanceScheduleAttribute{}
	for _, item := range api.Schedule {
		s := CheckMaintenanceScheduleAttribute{
			Type:         types.StringValue(item.Type),
			Monthday:     types.Int32Value(int32(item.Monthday)),
			MonthdayFrom: types.Int32Value(int32(item.MonthdayFrom)),
			MonthdayTo:   types.Int32Value(int32(item.MonthdayTo)),
		}

		if item.Type != "ONCE" {
			s.FromTime = types.StringValue(item.FromTime)
			s.ToTime = types.StringValue(item.ToTime)
		} else {
			s.FromTime = types.StringNull()
			s.ToTime = types.StringNull()
		}

		if len(item.Weekdays) > 0 {
			wd := make([]int32, len(item.Weekdays))
			for i := range item.Weekdays {
				wd[i] = int32(item.Weekdays[i])
			}
			s.Weekdays = a.SliceValue(wd)
		} else {
			s.Weekdays = types.SetNull(types.Int32Type)
		}

		var d diag.Diagnostics
		if item.OnceStartDate != "" {
			s.OnceStartDate, d = timetypes.NewRFC3339PointerValue(&item.OnceStartDate)
			if d.HasError() {
				return nil, fmt.Errorf("error parsing OnceStartDate: %v", d)
			}
		}

		if item.OnceEndDate != "" {
			s.OnceEndDate, d = timetypes.NewRFC3339PointerValue(&item.OnceEndDate)
			if d.HasError() {
				return nil, fmt.Errorf("error parsing OnceEndDate: %v", d)
			}
		}

		if item.Monthday != 0 {
			s.MonthdayTo = types.Int32Value(int32(0))
			s.MonthdayFrom = types.Int32Value(int32(0))
		}
		schedule = append(schedule, s)
	}
	var diags diag.Diagnostics
	if model.Schedule, diags = a.ScheduleValue(schedule); diags.HasError() {
		return nil, fmt.Errorf("failed to convert schedule: %v", diags)
	}

	if len(schedule) > 0 {
		model.schedule = schedule
	}
	return &model, nil
}

type CheckMaintenanceResourceAPI struct {
	provider *providerImpl
}

func (a CheckMaintenanceResourceAPI) Create(ctx context.Context, arg CheckMaintenanceWrapper) (*CheckMaintenanceWrapper, error) {
	pk := upapi.PrimaryKey(arg.CheckID)
	return a.Update(ctx, pk, arg)
}

func (a CheckMaintenanceResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*CheckMaintenanceWrapper, error) {
	obj, err := a.provider.api.Checks().Get(ctx, arg)
	if err != nil {
		return nil, err
	}
	var m upapi.CheckMaintenance
	if obj.Maintenance != nil {
		m = *obj.Maintenance
	}
	return &CheckMaintenanceWrapper{CheckMaintenance: m, CheckID: obj.PK}, nil
}

func (a CheckMaintenanceResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg CheckMaintenanceWrapper) (*CheckMaintenanceWrapper, error) {
	obj, err := a.provider.api.Checks().UpdateMaintenance(ctx, pk, arg.CheckMaintenance)
	if err != nil {
		return nil, err
	}
	var m upapi.CheckMaintenance
	if obj.Maintenance != nil {
		m = *obj.Maintenance
	}
	return &CheckMaintenanceWrapper{CheckMaintenance: m, CheckID: obj.PK}, nil
}

func (a CheckMaintenanceResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	_, err := a.provider.api.Checks().UpdateMaintenance(ctx, pk, upapi.CheckMaintenance{State: "ACTIVE"})
	return err
}
