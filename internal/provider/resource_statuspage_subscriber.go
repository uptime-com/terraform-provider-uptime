package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageSubscriberResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[StatusPageSubscriberResourceModel, StatusPageSubscriberWrapper, StatusPageSubscriberWrapper](
		&StatusPageSubscriberResourceAPI{provider: p},
		StatusPageSubscriberResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "statuspage_subscriber",
			Schema: schema.Schema{
				Description: "Status page subscriber resource. Import using composite ID: `terraform import uptime_statuspage_subscriber.example statuspage_id:subscriber_id`",
				Attributes: map[string]schema.Attribute{
					"statuspage_id": schema.Int64Attribute{
						Required: true,
					},
					"id": ComputedIDSchemaAttribute(), // Uses delete+create for updates
					"target": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							OneOfStringValidator([]string{"EMAIL", "SMS", "SLACK", "WEBHOOK"}),
						},
					},
					"force_validation_sms": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
				},
			},
		},
		ImportStateCompositeID,
	)
}

type StatusPageSubscriberWrapper struct {
	upapi.StatusPageSubscriber

	StatusPageID int64
}

func (w StatusPageSubscriberWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.PK)
}

type StatusPageSubscriberResourceModel struct {
	StatusPageID       types.Int64  `tfsdk:"statuspage_id"`
	ID                 types.Int64  `tfsdk:"id"`
	Target             types.String `tfsdk:"target"`
	Type               types.String `tfsdk:"type"`
	ForceValidationSMS types.Bool   `tfsdk:"force_validation_sms"`
}

func (m StatusPageSubscriberResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageSubscriberResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a StatusPageSubscriberResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*StatusPageSubscriberResourceModel, diag.Diagnostics) {
	var model StatusPageSubscriberResourceModel
	if diags := sg.Get(ctx, &model); diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a StatusPageSubscriberResourceModelAdapter) ToAPIArgument(
	model StatusPageSubscriberResourceModel,
) (*StatusPageSubscriberWrapper, error) {
	return &StatusPageSubscriberWrapper{
		StatusPageID: model.StatusPageID.ValueInt64(),
		StatusPageSubscriber: upapi.StatusPageSubscriber{
			PK:                 model.ID.ValueInt64(),
			Target:             model.Target.ValueString(),
			Type:               model.Type.ValueString(),
			ForceValidationSMS: model.ForceValidationSMS.ValueBool(),
		},
	}, nil
}

func (a StatusPageSubscriberResourceModelAdapter) FromAPIResult(
	api StatusPageSubscriberWrapper,
) (*StatusPageSubscriberResourceModel, error) {
	return &StatusPageSubscriberResourceModel{
		StatusPageID:       types.Int64Value(api.StatusPageID),
		ID:                 types.Int64Value(api.PK),
		Target:             types.StringValue(api.Target),
		Type:               types.StringValue(api.Type),
		ForceValidationSMS: types.BoolValue(api.ForceValidationSMS),
	}, nil
}

type StatusPageSubscriberResourceAPI struct {
	provider *providerImpl
}

func (a StatusPageSubscriberResourceAPI) Create(ctx context.Context, arg StatusPageSubscriberWrapper) (*StatusPageSubscriberWrapper, error) {
	obj, err := a.provider.api.StatusPages().Subscribers(upapi.PrimaryKey(arg.StatusPageID)).Create(ctx, arg.StatusPageSubscriber)
	if err != nil {
		return nil, err
	}

	return &StatusPageSubscriberWrapper{StatusPageSubscriber: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageSubscriberResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*StatusPageSubscriberWrapper, error) {
	model, ok := arg.(StatusPageSubscriberResourceModel)
	if !ok {
		return nil, fmt.Errorf("resource read failed:unexpected type %T", arg)
	}

	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	obj, err := a.provider.api.StatusPages().Subscribers(statusPageID).Get(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &StatusPageSubscriberWrapper{StatusPageSubscriber: *obj, StatusPageID: int64(statusPageID)}, nil
}

func (a StatusPageSubscriberResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg StatusPageSubscriberWrapper) (*StatusPageSubscriberWrapper, error) {
	err := a.provider.api.StatusPages().Subscribers(upapi.PrimaryKey(arg.StatusPageID)).Delete(ctx, pk)
	if err != nil {
		return nil, err
	}

	obj, err := a.provider.api.StatusPages().Subscribers(upapi.PrimaryKey(arg.StatusPageID)).Create(ctx, arg.StatusPageSubscriber)
	if err != nil {
		return nil, err
	}
	return &StatusPageSubscriberWrapper{StatusPageSubscriber: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageSubscriberResourceAPI) Delete(ctx context.Context, arg upapi.PrimaryKeyable) error {
	model, ok := arg.(StatusPageSubscriberResourceModel)
	if !ok {
		return fmt.Errorf("resource delete failed: unexpected type %T", arg)
	}
	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	return a.provider.api.StatusPages().Subscribers(statusPageID).Delete(ctx, arg)
}
