package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewMaintenanceScheduleResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[MaintenanceScheduleModel, upapi.MaintenanceScheduleInput, upapi.MaintenanceSchedule](
		&MaintenanceScheduleAPI{provider: p},
		MaintenanceScheduleModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "maintenance_schedule",
			Schema: schema.Schema{
				Description: "Account-level maintenance schedule (maintenance window) targeting checks by service ID or tag ID. schedule_type RRULE/ONE_OFF only. Note: delete is a soft-delete server-side; deletion outside Terraform is not detected as drift.",
				Attributes: map[string]schema.Attribute{
					// IDSchemaAttribute() sets Computed + int64planmodifier.UseStateForUnknown()
					// (see attr.go). Every resource in this repo uses it for a stable id.
					"id": IDSchemaAttribute(),
					"name": schema.StringAttribute{
						Required: true,
					},
					"schedule_type": schema.StringAttribute{
						Required: true,
						Description: "Recurrence type. `RRULE` requires `rrule` and `duration_minutes`; " +
							"`ONE_OFF` requires `duration_minutes` or `ends_at`.",
						Validators: []validator.String{
							OneOfStringValidator([]string{"RRULE", "ONE_OFF"}),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"starts_at": schema.StringAttribute{
						Required:   true,
						CustomType: timetypes.RFC3339Type{},
					},
					"ends_at": schema.StringAttribute{
						Optional:   true,
						Computed:   true,
						CustomType: timetypes.RFC3339Type{},
						Description: "End time (RFC 3339). For ONE_OFF, computed from " +
							"starts_at + duration_minutes if omitted.",
					},
					"rrule": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Description: "RFC 5545 recurrence rule string (e.g. `FREQ=WEEKLY;BYDAY=SA`). " +
							"Required when schedule_type is RRULE.",
					},
					"duration_minutes": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Description: "Maintenance window length in minutes. Required for RRULE; " +
							"for ONE_OFF, provide this or ends_at.",
					},
					"is_active": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"pause_checks_during_maintenance": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					// Optional+Computed: FromAPIResult always returns a (possibly empty)
					// known set; without Computed, an omitted (null) config vs a returned
					// empty set causes "inconsistent result after apply".
					"services": schema.SetAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.Int64Type,
						Description: "Service (check) IDs this maintenance applies to. Use `uptime_check_*.id`.",
					},
					"tags": schema.SetAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.Int64Type,
						Description: "Service tag IDs this maintenance targets. Use `uptime_tag.id`.",
					},
					"created_at": schema.StringAttribute{
						Computed:   true,
						CustomType: timetypes.RFC3339Type{},
					},
					"modified_at": schema.StringAttribute{
						Computed:   true,
						CustomType: timetypes.RFC3339Type{},
					},
				},
			},
			ConfigValidators: func(_ context.Context) []resource.ConfigValidator {
				return []resource.ConfigValidator{maintenanceScheduleConfigValidator{}}
			},
		},
		ImportStateSimpleID,
	)
}

type MaintenanceScheduleModel struct {
	ID                           types.Int64       `tfsdk:"id"`
	Name                         types.String      `tfsdk:"name"`
	ScheduleType                 types.String      `tfsdk:"schedule_type"`
	StartsAt                     timetypes.RFC3339 `tfsdk:"starts_at"`
	EndsAt                       timetypes.RFC3339 `tfsdk:"ends_at"`
	RRule                        types.String      `tfsdk:"rrule"`
	DurationMinutes              types.Int64       `tfsdk:"duration_minutes"`
	IsActive                     types.Bool        `tfsdk:"is_active"`
	PauseChecksDuringMaintenance types.Bool        `tfsdk:"pause_checks_during_maintenance"`
	Services                     types.Set         `tfsdk:"services"`
	Tags                         types.Set         `tfsdk:"tags"`
	CreatedAt                    timetypes.RFC3339 `tfsdk:"created_at"`
	ModifiedAt                   timetypes.RFC3339 `tfsdk:"modified_at"`
}

func (m MaintenanceScheduleModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type MaintenanceScheduleModelAdapter struct {
	SetAttributeAdapter[int64]
}

func (a MaintenanceScheduleModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*MaintenanceScheduleModel, diag.Diagnostics) {
	model := *new(MaintenanceScheduleModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a MaintenanceScheduleModelAdapter) ToAPIArgument(
	model MaintenanceScheduleModel,
) (*upapi.MaintenanceScheduleInput, error) {
	arg := upapi.MaintenanceScheduleInput{
		Name:                         model.Name.ValueString(),
		ScheduleType:                 model.ScheduleType.ValueString(),
		StartsAt:                     model.StartsAt.ValueString(),
		RRule:                        model.RRule.ValueString(),
		IsActive:                     model.IsActive.ValueBool(),
		PauseChecksDuringMaintenance: model.PauseChecksDuringMaintenance.ValueBool(),
		Services:                     a.Slice(model.Services),
		Tags:                         a.Slice(model.Tags),
	}
	// Guarantee non-nil slices so empty targets serialize as `[]` not `null`
	// (the API rejects null here).
	if arg.Services == nil {
		arg.Services = []int64{}
	}
	if arg.Tags == nil {
		arg.Tags = []int64{}
	}
	if !model.EndsAt.IsNull() && !model.EndsAt.IsUnknown() {
		v := model.EndsAt.ValueString()
		arg.EndsAt = &v
	}
	if !model.DurationMinutes.IsNull() && !model.DurationMinutes.IsUnknown() {
		v := model.DurationMinutes.ValueInt64()
		arg.DurationMinutes = &v
	}
	return &arg, nil
}

func (a MaintenanceScheduleModelAdapter) FromAPIResult(
	api upapi.MaintenanceSchedule,
) (*MaintenanceScheduleModel, error) {
	model := MaintenanceScheduleModel{
		ID:                           types.Int64Value(api.PK),
		Name:                         types.StringValue(api.Name),
		ScheduleType:                 types.StringValue(api.ScheduleType),
		IsActive:                     types.BoolValue(api.IsActive),
		PauseChecksDuringMaintenance: types.BoolValue(api.PauseChecksDuringMaintenance),
	}
	// rrule is Optional+Computed: map empty -> Null to avoid ""-vs-null drift on
	// ONE_OFF schedules (which have no rrule).
	if api.RRule != "" {
		model.RRule = types.StringValue(api.RRule)
	} else {
		model.RRule = types.StringNull()
	}
	var err error
	if model.StartsAt, err = parseRFC3339OrNull(api.StartsAt); err != nil {
		return nil, err
	}
	if model.EndsAt, err = parseRFC3339OrNull(api.EndsAt); err != nil {
		return nil, err
	}
	if model.CreatedAt, err = parseRFC3339OrNull(api.CreatedAt); err != nil {
		return nil, err
	}
	if model.ModifiedAt, err = parseRFC3339OrNull(api.ModifiedAt); err != nil {
		return nil, err
	}
	if api.DurationMinutes != nil {
		model.DurationMinutes = types.Int64Value(*api.DurationMinutes)
	} else {
		model.DurationMinutes = types.Int64Null()
	}
	serviceIDs := make([]int64, len(api.Services))
	for i := range api.Services {
		serviceIDs[i] = api.Services[i].PK
	}
	model.Services = a.SliceValue(serviceIDs)
	tagIDs := make([]int64, len(api.Tags))
	for i := range api.Tags {
		tagIDs[i] = api.Tags[i].PK
	}
	model.Tags = a.SliceValue(tagIDs)
	return &model, nil
}

type MaintenanceScheduleAPI struct {
	provider *providerImpl
}

func (a MaintenanceScheduleAPI) Create(
	ctx context.Context, arg upapi.MaintenanceScheduleInput,
) (*upapi.MaintenanceSchedule, error) {
	return a.provider.api.MaintenanceSchedules().Create(ctx, arg)
}

func (a MaintenanceScheduleAPI) Read(
	ctx context.Context, pk upapi.PrimaryKeyable,
) (*upapi.MaintenanceSchedule, error) {
	return a.provider.api.MaintenanceSchedules().Get(ctx, pk)
}

func (a MaintenanceScheduleAPI) Update(
	ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.MaintenanceScheduleInput,
) (*upapi.MaintenanceSchedule, error) {
	return a.provider.api.MaintenanceSchedules().Update(ctx, pk, arg)
}

func (a MaintenanceScheduleAPI) Delete(
	ctx context.Context, pk upapi.PrimaryKeyable,
) error {
	return a.provider.api.MaintenanceSchedules().Delete(ctx, pk)
}

// parseRFC3339OrNull converts an API timestamp string into a timetypes.RFC3339.
// It returns Null for empty strings and an error for malformed timestamps. This
// shared helper is reused by the maintenance notification resource.
func parseRFC3339OrNull(s string) (timetypes.RFC3339, error) {
	if s == "" {
		return timetypes.NewRFC3339Null(), nil
	}
	v, diags := timetypes.NewRFC3339Value(s)
	if diags.HasError() {
		return timetypes.RFC3339{}, fmt.Errorf("invalid RFC3339 %q: %v", s, diags)
	}
	return v, nil
}

type maintenanceScheduleConfigValidator struct{}

func (v maintenanceScheduleConfigValidator) Description(_ context.Context) string {
	return "validates schedule_type field requirements"
}

func (v maintenanceScheduleConfigValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v maintenanceScheduleConfigValidator) ValidateResource(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var model MaintenanceScheduleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() || model.ScheduleType.IsUnknown() || model.ScheduleType.IsNull() {
		return
	}
	switch model.ScheduleType.ValueString() {
	case "RRULE":
		if model.RRule.IsNull() || model.RRule.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(path.Root("rrule"), "Missing rrule",
				"schedule_type RRULE requires a non-empty rrule.")
		}
		if model.DurationMinutes.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("duration_minutes"),
				"Missing duration_minutes", "schedule_type RRULE requires duration_minutes.")
		}
	case "ONE_OFF":
		if model.DurationMinutes.IsNull() && model.EndsAt.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("duration_minutes"),
				"Missing duration_minutes or ends_at",
				"schedule_type ONE_OFF requires duration_minutes or ends_at.")
		}
	}
}
