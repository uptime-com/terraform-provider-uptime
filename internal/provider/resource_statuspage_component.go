package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageComponentResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[StatusPageComponentResourceModel, StatusPageComponentWrapper, StatusPageComponentWrapper]{
		api: &StatusPageComponentResourceAPI{provider: p},
		mod: StatusPageComponentResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "statuspage_component",
			Schema: schema.Schema{
				Description: "Status page component resource",
				Attributes: map[string]schema.Attribute{
					"statuspage_id": schema.Int64Attribute{
						Required: true,
					},
					"id":   IDSchemaAttribute(),
					"url":  URLSchemaAttribute(),
					"name": NameSchemaAttribute(),
					"description": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Default:  stringdefault.StaticString(""),
					},
					"group_id": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"service_id": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"is_group": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"status": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Validators: []validator.String{
							OneOfStringValidator([]string{"operational", "major-outage", "partial-outage", "degraded-performance", "under-maintenance"}),
						},
					},
					"auto_status_down": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Validators: []validator.String{
							OneOfStringValidator([]string{"major-outage", "partial-outage", "degraded-performance", "under-maintenance"}),
						},
					},
					"auto_status_up": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Validators: []validator.String{
							OneOfStringValidator([]string{"operational", "major-outage", "partial-outage", "degraded-performance", "under-maintenance"}),
						},
					},
				},
			},
		},
	}
}

type StatusPageComponentWrapper struct {
	upapi.StatusPageComponent

	StatusPageID int64
}

func (w StatusPageComponentWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.PK)
}

type StatusPageComponentResourceModel struct {
	StatusPageID   types.Int64  `tfsdk:"statuspage_id"`
	ID             types.Int64  `tfsdk:"id"`
	URL            types.String `tfsdk:"url"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	GroupID        types.Int64  `tfsdk:"group_id"`
	ServiceID      types.Int64  `tfsdk:"service_id"`
	IsGroup        types.Bool   `tfsdk:"is_group"`
	Status         types.String `tfsdk:"status"`
	AutoStatusDown types.String `tfsdk:"auto_status_down"`
	AutoStatusUp   types.String `tfsdk:"auto_status_up"`
}

func (m StatusPageComponentResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageComponentResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a StatusPageComponentResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*StatusPageComponentResourceModel, diag.Diagnostics) {
	var model StatusPageComponentResourceModel
	if diags := sg.Get(ctx, &model); diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a StatusPageComponentResourceModelAdapter) ToAPIArgument(
	model StatusPageComponentResourceModel,
) (*StatusPageComponentWrapper, error) {
	return &StatusPageComponentWrapper{
		StatusPageID: model.StatusPageID.ValueInt64(),
		StatusPageComponent: upapi.StatusPageComponent{
			Name:           model.Name.ValueString(),
			Description:    model.Description.ValueString(),
			IsGroup:        model.IsGroup.ValueBool(),
			GroupID:        model.GroupID.ValueInt64(),
			ServiceID:      model.ServiceID.ValueInt64(),
			Status:         model.Status.ValueString(),
			AutoStatusDown: model.AutoStatusDown.ValueString(),
			AutoStatusUp:   model.AutoStatusUp.ValueString(),
		},
	}, nil
}

func (a StatusPageComponentResourceModelAdapter) FromAPIResult(
	api StatusPageComponentWrapper,
) (*StatusPageComponentResourceModel, error) {
	return &StatusPageComponentResourceModel{
		StatusPageID:   types.Int64Value(api.StatusPageID),
		ID:             types.Int64Value(api.PK),
		URL:            types.StringValue(api.URL),
		Name:           types.StringValue(api.Name),
		Description:    types.StringValue(api.Description),
		IsGroup:        types.BoolValue(api.IsGroup),
		GroupID:        types.Int64Value(api.GroupID),
		ServiceID:      types.Int64Value(api.ServiceID),
		Status:         types.StringValue(api.Status),
		AutoStatusDown: types.StringValue(api.AutoStatusDown),
		AutoStatusUp:   types.StringValue(api.AutoStatusUp),
	}, nil
}

type StatusPageComponentResourceAPI struct {
	provider *providerImpl
}

func (a StatusPageComponentResourceAPI) Create(ctx context.Context, arg StatusPageComponentWrapper) (*StatusPageComponentWrapper, error) {
	obj, err := a.provider.api.StatusPages().Components(upapi.PrimaryKey(arg.StatusPageID)).Create(ctx, arg.StatusPageComponent)
	if err != nil {
		return nil, err
	}

	return &StatusPageComponentWrapper{StatusPageComponent: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageComponentResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*StatusPageComponentWrapper, error) {
	model, ok := arg.(StatusPageComponentResourceModel)
	if !ok {
		return nil, fmt.Errorf("resource read failed:unexpected type %T", arg)
	}

	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	obj, err := a.provider.api.StatusPages().Components(statusPageID).Get(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &StatusPageComponentWrapper{StatusPageComponent: *obj, StatusPageID: int64(statusPageID)}, nil
}

func (a StatusPageComponentResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg StatusPageComponentWrapper) (*StatusPageComponentWrapper, error) {
	obj, err := a.provider.api.StatusPages().Components(upapi.PrimaryKey(arg.StatusPageID)).Update(ctx, pk, arg.StatusPageComponent)
	if err != nil {
		return nil, err
	}
	return &StatusPageComponentWrapper{StatusPageComponent: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageComponentResourceAPI) Delete(ctx context.Context, arg upapi.PrimaryKeyable) error {
	model, ok := arg.(StatusPageComponentResourceModel)
	if !ok {
		return fmt.Errorf("resource delete failed: unexpected type %T", arg)
	}
	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	return a.provider.api.StatusPages().Components(statusPageID).Delete(ctx, arg)
}
