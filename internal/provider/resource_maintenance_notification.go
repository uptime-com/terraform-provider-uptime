package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewMaintenanceNotificationResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[MaintenanceNotificationModel, upapi.MaintenanceScheduleNotificationInput, upapi.MaintenanceScheduleNotification](
		&MaintenanceNotificationAPI{provider: p},
		MaintenanceNotificationModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "maintenance_notification",
			Schema: schema.Schema{
				Description: "Notification rule for a maintenance schedule (alert N seconds before/after START or END).",
				Attributes: map[string]schema.Attribute{
					"id": IDSchemaAttribute(),
					"schedule_id": schema.Int64Attribute{
						Required: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"offset": schema.Int64Attribute{
						Required:    true,
						Description: "Offset in seconds relative to the event; negative means before.",
					},
					"event": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							OneOfStringValidator([]string{"START", "END"}),
						},
					},
					"contact_groups": schema.SetAttribute{
						Required:    true,
						ElementType: types.Int64Type,
						Description: "Contact group IDs to notify. Use `uptime_contact.<name>.id`.",
					},
					"created_at":  schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
					"modified_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type MaintenanceNotificationModel struct {
	ID            types.Int64       `tfsdk:"id"`
	ScheduleID    types.Int64       `tfsdk:"schedule_id"`
	Offset        types.Int64       `tfsdk:"offset"`
	Event         types.String      `tfsdk:"event"`
	ContactGroups types.Set         `tfsdk:"contact_groups"`
	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at"`
	ModifiedAt    timetypes.RFC3339 `tfsdk:"modified_at"`
}

func (m MaintenanceNotificationModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type MaintenanceNotificationModelAdapter struct {
	SetAttributeAdapter[int64]
}

func (a MaintenanceNotificationModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*MaintenanceNotificationModel, diag.Diagnostics) {
	model := *new(MaintenanceNotificationModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a MaintenanceNotificationModelAdapter) ToAPIArgument(
	model MaintenanceNotificationModel,
) (*upapi.MaintenanceScheduleNotificationInput, error) {
	// contact_groups must serialize as `[]` not `null` (the API rejects null here).
	contactGroups := a.Slice(model.ContactGroups)
	if contactGroups == nil {
		contactGroups = []int64{}
	}
	return &upapi.MaintenanceScheduleNotificationInput{
		ScheduleID:    model.ScheduleID.ValueInt64(),
		Offset:        model.Offset.ValueInt64(),
		Event:         model.Event.ValueString(),
		ContactGroups: contactGroups,
	}, nil
}

func (a MaintenanceNotificationModelAdapter) FromAPIResult(
	api upapi.MaintenanceScheduleNotification,
) (*MaintenanceNotificationModel, error) {
	model := MaintenanceNotificationModel{
		ID:         types.Int64Value(api.PK),
		ScheduleID: types.Int64Value(api.ScheduleID),
		Offset:     types.Int64Value(api.Offset),
		Event:      types.StringValue(api.Event),
	}
	ids := make([]int64, len(api.ContactGroups))
	for i := range api.ContactGroups {
		ids[i] = api.ContactGroups[i].PK
	}
	model.ContactGroups = a.SliceValue(ids)
	createdAt, err := parseRFC3339OrNull(api.CreatedAt)
	if err != nil {
		return nil, err
	}
	model.CreatedAt = createdAt
	modifiedAt, err := parseRFC3339OrNull(api.ModifiedAt)
	if err != nil {
		return nil, err
	}
	model.ModifiedAt = modifiedAt
	return &model, nil
}

type MaintenanceNotificationAPI struct {
	provider *providerImpl
}

func (a MaintenanceNotificationAPI) Create(
	ctx context.Context, arg upapi.MaintenanceScheduleNotificationInput,
) (*upapi.MaintenanceScheduleNotification, error) {
	return a.provider.api.MaintenanceNotifications().Create(ctx, arg)
}

func (a MaintenanceNotificationAPI) Read(
	ctx context.Context, pk upapi.PrimaryKeyable,
) (*upapi.MaintenanceScheduleNotification, error) {
	return a.provider.api.MaintenanceNotifications().Get(ctx, pk)
}

func (a MaintenanceNotificationAPI) Update(
	ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.MaintenanceScheduleNotificationInput,
) (*upapi.MaintenanceScheduleNotification, error) {
	return a.provider.api.MaintenanceNotifications().Update(ctx, pk, arg)
}

func (a MaintenanceNotificationAPI) Delete(
	ctx context.Context, pk upapi.PrimaryKeyable,
) error {
	return a.provider.api.MaintenanceNotifications().Delete(ctx, pk)
}
