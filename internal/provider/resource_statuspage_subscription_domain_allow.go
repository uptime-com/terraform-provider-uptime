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

func NewStatusPageSubsDomainAllowResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[StatusPageSubsDomainAllowResourceModel, StatusPageSubsDomainAllowWrapper, StatusPageSubsDomainAllowWrapper]{
		api: &StatusPageSubsDomainAllowResourceAPI{provider: p},
		mod: StatusPageSubsDomainAllowResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "statuspage_subscription_domain_allow",
			Schema: schema.Schema{
				Description: "Status page subscription domain allow resource",
				Attributes: map[string]schema.Attribute{
					"statuspage_id": schema.Int64Attribute{
						Required: true,
					},
					"id":     IDSchemaAttribute(),
					"domain": NameSchemaAttribute(),
				},
			},
		},
	}
}

type StatusPageSubsDomainAllowWrapper struct {
	upapi.StatusPageSubsDomainAllowList

	StatusPageID int64
}

func (w StatusPageSubsDomainAllowWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.PK)
}

type StatusPageSubsDomainAllowResourceModel struct {
	StatusPageID types.Int64  `tfsdk:"statuspage_id"`
	ID           types.Int64  `tfsdk:"id"`
	Domain       types.String `tfsdk:"domain"`
}

func (m StatusPageSubsDomainAllowResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageSubsDomainAllowResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a StatusPageSubsDomainAllowResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*StatusPageSubsDomainAllowResourceModel, diag.Diagnostics) {
	var model StatusPageSubsDomainAllowResourceModel
	if diags := sg.Get(ctx, &model); diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a StatusPageSubsDomainAllowResourceModelAdapter) ToAPIArgument(
	model StatusPageSubsDomainAllowResourceModel,
) (*StatusPageSubsDomainAllowWrapper, error) {
	return &StatusPageSubsDomainAllowWrapper{
		StatusPageID: model.StatusPageID.ValueInt64(),
		StatusPageSubsDomainAllowList: upapi.StatusPageSubsDomainAllowList{
			PK:     model.ID.ValueInt64(),
			Domain: model.Domain.ValueString(),
		},
	}, nil
}

func (a StatusPageSubsDomainAllowResourceModelAdapter) FromAPIResult(
	api StatusPageSubsDomainAllowWrapper,
) (*StatusPageSubsDomainAllowResourceModel, error) {
	return &StatusPageSubsDomainAllowResourceModel{
		StatusPageID: types.Int64Value(api.StatusPageID),
		ID:           types.Int64Value(api.PK),
		Domain:       types.StringValue(api.Domain),
	}, nil
}

type StatusPageSubsDomainAllowResourceAPI struct {
	provider *providerImpl
}

func (a StatusPageSubsDomainAllowResourceAPI) Create(
	ctx context.Context, arg StatusPageSubsDomainAllowWrapper,
) (*StatusPageSubsDomainAllowWrapper, error) {
	obj, err := a.provider.api.StatusPages().
		SubscriptionDomainAllowList(upapi.PrimaryKey(arg.StatusPageID)).
		Create(ctx, arg.StatusPageSubsDomainAllowList)
	if err != nil {
		return nil, err
	}

	return &StatusPageSubsDomainAllowWrapper{
		StatusPageSubsDomainAllowList: *obj,
		StatusPageID:                  arg.StatusPageID,
	}, nil
}

func (a StatusPageSubsDomainAllowResourceAPI) Read(
	ctx context.Context, arg upapi.PrimaryKeyable,
) (*StatusPageSubsDomainAllowWrapper, error) {
	model, ok := arg.(StatusPageSubsDomainAllowResourceModel)
	if !ok {
		return nil, fmt.Errorf("resource read failed:unexpected type %T", arg)
	}

	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	obj, err := a.provider.api.StatusPages().SubscriptionDomainAllowList(statusPageID).Get(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &StatusPageSubsDomainAllowWrapper{
		StatusPageSubsDomainAllowList: *obj,
		StatusPageID:                  int64(statusPageID),
	}, nil
}

func (a StatusPageSubsDomainAllowResourceAPI) Update(
	ctx context.Context, pk upapi.PrimaryKeyable, arg StatusPageSubsDomainAllowWrapper,
) (*StatusPageSubsDomainAllowWrapper, error) {
	obj, err := a.provider.api.StatusPages().
		SubscriptionDomainAllowList(upapi.PrimaryKey(arg.StatusPageID)).
		Update(ctx, pk, arg.StatusPageSubsDomainAllowList)
	if err != nil {
		return nil, err
	}
	return &StatusPageSubsDomainAllowWrapper{
		StatusPageSubsDomainAllowList: *obj,
		StatusPageID:                  arg.StatusPageID,
	}, nil
}

func (a StatusPageSubsDomainAllowResourceAPI) Delete(ctx context.Context, arg upapi.PrimaryKeyable) error {
	model, ok := arg.(StatusPageSubsDomainAllowResourceModel)
	if !ok {
		return fmt.Errorf("resource delete failed: unexpected type %T", arg)
	}
	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	return a.provider.api.StatusPages().SubscriptionDomainAllowList(statusPageID).Delete(ctx, arg)
}
