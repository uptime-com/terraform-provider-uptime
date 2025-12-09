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

func NewStatusPageSubsDomainBlockResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[StatusPageSubsDomainBlockResourceModel, StatusPageSubsDomainBlockWrapper, StatusPageSubsDomainBlockWrapper](
		&StatusPageSubsDomainBlockResourceAPI{provider: p},
		StatusPageSubsDomainBlockResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "statuspage_subscription_domain_block",
			Schema: schema.Schema{
				Description: "Status page subscription domain block resource. Import using composite ID: `terraform import uptime_statuspage_subscription_domain_block.example statuspage_id:domain_id`",
				Attributes: map[string]schema.Attribute{
					"statuspage_id": schema.Int64Attribute{
						Required: true,
					},
					"id":     IDSchemaAttribute(),
					"domain": NameSchemaAttribute(),
				},
			},
		},
		ImportStateCompositeID,
	)
}

type StatusPageSubsDomainBlockWrapper struct {
	upapi.StatusPageSubsDomainBlockList

	StatusPageID int64
}

func (w StatusPageSubsDomainBlockWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.PK)
}

type StatusPageSubsDomainBlockResourceModel struct {
	StatusPageID types.Int64  `tfsdk:"statuspage_id"`
	ID           types.Int64  `tfsdk:"id"`
	Domain       types.String `tfsdk:"domain"`
}

func (m StatusPageSubsDomainBlockResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageSubsDomainBlockResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a StatusPageSubsDomainBlockResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*StatusPageSubsDomainBlockResourceModel, diag.Diagnostics) {
	var model StatusPageSubsDomainBlockResourceModel
	if diags := sg.Get(ctx, &model); diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a StatusPageSubsDomainBlockResourceModelAdapter) ToAPIArgument(
	model StatusPageSubsDomainBlockResourceModel,
) (*StatusPageSubsDomainBlockWrapper, error) {
	return &StatusPageSubsDomainBlockWrapper{
		StatusPageID: model.StatusPageID.ValueInt64(),
		StatusPageSubsDomainBlockList: upapi.StatusPageSubsDomainBlockList{
			PK:     model.ID.ValueInt64(),
			Domain: model.Domain.ValueString(),
		},
	}, nil
}

func (a StatusPageSubsDomainBlockResourceModelAdapter) FromAPIResult(
	api StatusPageSubsDomainBlockWrapper,
) (*StatusPageSubsDomainBlockResourceModel, error) {
	return &StatusPageSubsDomainBlockResourceModel{
		StatusPageID: types.Int64Value(api.StatusPageID),
		ID:           types.Int64Value(api.PK),
		Domain:       types.StringValue(api.Domain),
	}, nil
}

type StatusPageSubsDomainBlockResourceAPI struct {
	provider *providerImpl
}

func (a StatusPageSubsDomainBlockResourceAPI) Create(
	ctx context.Context, arg StatusPageSubsDomainBlockWrapper,
) (*StatusPageSubsDomainBlockWrapper, error) {
	obj, err := a.provider.api.StatusPages().
		SubscriptionDomainBlockList(upapi.PrimaryKey(arg.StatusPageID)).
		Create(ctx, arg.StatusPageSubsDomainBlockList)
	if err != nil {
		return nil, err
	}

	return &StatusPageSubsDomainBlockWrapper{
		StatusPageSubsDomainBlockList: *obj,
		StatusPageID:                  arg.StatusPageID,
	}, nil
}

func (a StatusPageSubsDomainBlockResourceAPI) Read(
	ctx context.Context, arg upapi.PrimaryKeyable,
) (*StatusPageSubsDomainBlockWrapper, error) {
	model, ok := arg.(StatusPageSubsDomainBlockResourceModel)
	if !ok {
		return nil, fmt.Errorf("resource read failed:unexpected type %T", arg)
	}

	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	obj, err := a.provider.api.StatusPages().SubscriptionDomainBlockList(statusPageID).Get(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &StatusPageSubsDomainBlockWrapper{
		StatusPageSubsDomainBlockList: *obj,
		StatusPageID:                  int64(statusPageID),
	}, nil
}

func (a StatusPageSubsDomainBlockResourceAPI) Update(
	ctx context.Context, pk upapi.PrimaryKeyable, arg StatusPageSubsDomainBlockWrapper,
) (*StatusPageSubsDomainBlockWrapper, error) {
	obj, err := a.provider.api.StatusPages().
		SubscriptionDomainBlockList(upapi.PrimaryKey(arg.StatusPageID)).
		Update(ctx, pk, arg.StatusPageSubsDomainBlockList)
	if err != nil {
		return nil, err
	}
	return &StatusPageSubsDomainBlockWrapper{
		StatusPageSubsDomainBlockList: *obj,
		StatusPageID:                  arg.StatusPageID,
	}, nil
}

func (a StatusPageSubsDomainBlockResourceAPI) Delete(ctx context.Context, arg upapi.PrimaryKeyable) error {
	model, ok := arg.(StatusPageSubsDomainBlockResourceModel)
	if !ok {
		return fmt.Errorf("resource delete failed: unexpected type %T", arg)
	}
	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	return a.provider.api.StatusPages().SubscriptionDomainBlockList(statusPageID).Delete(ctx, arg)
}
