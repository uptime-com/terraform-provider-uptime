package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageUserResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[StatusPageUserResourceModel, StatusPageUserWrapper, StatusPageUserWrapper](
		&StatusPageUserResourceAPI{provider: p},
		StatusPageUserResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "statuspage_user",
			Schema: schema.Schema{
				Description: "Status page user resource. Import using composite ID: `terraform import uptime_statuspage_user.example statuspage_id:user_id`",
				Attributes: map[string]schema.Attribute{
					"statuspage_id": schema.Int64Attribute{
						Required: true,
					},
					"id": IDSchemaAttribute(),
					"email": schema.StringAttribute{
						Required: true,
					},
					"first_name": schema.StringAttribute{
						Required: true,
					},
					"last_name": schema.StringAttribute{
						Required: true,
					},
					"is_active": schema.BoolAttribute{
						Required: true,
					},
				},
			},
		},
		ImportStateCompositeID,
	)
}

type StatusPageUserWrapper struct {
	upapi.StatusPageUser

	StatusPageID int64
}

func (w StatusPageUserWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.PK)
}

type StatusPageUserResourceModel struct {
	StatusPageID types.Int64  `tfsdk:"statuspage_id"`
	ID           types.Int64  `tfsdk:"id"`
	Email        types.String `tfsdk:"email"`
	FirstName    types.String `tfsdk:"first_name"`
	LastName     types.String `tfsdk:"last_name"`
	IsActive     types.Bool   `tfsdk:"is_active"`
}

func (m StatusPageUserResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageUserResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a StatusPageUserResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*StatusPageUserResourceModel, diag.Diagnostics) {
	var model StatusPageUserResourceModel
	if diags := sg.Get(ctx, &model); diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a StatusPageUserResourceModelAdapter) ToAPIArgument(
	model StatusPageUserResourceModel,
) (*StatusPageUserWrapper, error) {
	return &StatusPageUserWrapper{
		StatusPageID: model.StatusPageID.ValueInt64(),
		StatusPageUser: upapi.StatusPageUser{
			PK:        model.ID.ValueInt64(),
			Email:     model.Email.ValueString(),
			FirstName: model.FirstName.ValueString(),
			LastName:  model.LastName.ValueString(),
			IsActive:  model.IsActive.ValueBool(),
		},
	}, nil
}

func (a StatusPageUserResourceModelAdapter) FromAPIResult(
	api StatusPageUserWrapper,
) (*StatusPageUserResourceModel, error) {
	return &StatusPageUserResourceModel{
		StatusPageID: types.Int64Value(api.StatusPageID),
		ID:           types.Int64Value(api.PK),
		Email:        types.StringValue(api.Email),
		FirstName:    types.StringValue(api.FirstName),
		LastName:     types.StringValue(api.LastName),
		IsActive:     types.BoolValue(api.IsActive),
	}, nil
}

type StatusPageUserResourceAPI struct {
	provider *providerImpl
}

func (a StatusPageUserResourceAPI) Create(ctx context.Context, arg StatusPageUserWrapper) (*StatusPageUserWrapper, error) {
	obj, err := a.provider.api.StatusPages().Users(upapi.PrimaryKey(arg.StatusPageID)).Create(ctx, arg.StatusPageUser)
	if err != nil {
		return nil, err
	}

	return &StatusPageUserWrapper{StatusPageUser: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageUserResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*StatusPageUserWrapper, error) {
	model, ok := arg.(StatusPageUserResourceModel)
	if !ok {
		return nil, fmt.Errorf("resource read failed:unexpected type %T", arg)
	}

	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	obj, err := a.provider.api.StatusPages().Users(statusPageID).Get(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &StatusPageUserWrapper{StatusPageUser: *obj, StatusPageID: int64(statusPageID)}, nil
}

func (a StatusPageUserResourceAPI) Update(
	ctx context.Context, pk upapi.PrimaryKeyable, arg StatusPageUserWrapper,
) (*StatusPageUserWrapper, error) {
	obj, err := a.provider.api.StatusPages().
		Users(upapi.PrimaryKey(arg.StatusPageID)).
		Update(ctx, pk, arg.StatusPageUser)
	if err != nil {
		return nil, err
	}
	return &StatusPageUserWrapper{StatusPageUser: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageUserResourceAPI) Delete(ctx context.Context, arg upapi.PrimaryKeyable) error {
	model, ok := arg.(StatusPageUserResourceModel)
	if !ok {
		return fmt.Errorf("resource delete failed: unexpected type %T", arg)
	}
	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	return a.provider.api.StatusPages().Users(statusPageID).Delete(ctx, arg)
}
