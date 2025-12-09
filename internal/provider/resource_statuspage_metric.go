package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageMetricResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[StatusPageMetricResourceModel, StatusPageMetricWrapper, StatusPageMetricWrapper](
		&StatusPageMetricResourceAPI{provider: p},
		StatusPageMetricResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "statuspage_metric",
			Schema: schema.Schema{
				Description: "Status page metric resource. Import using composite ID: `terraform import uptime_statuspage_metric.example statuspage_id:metric_id`",
				Attributes: map[string]schema.Attribute{
					"statuspage_id": schema.Int64Attribute{
						Required: true,
					},
					"id":   IDSchemaAttribute(),
					"url":  URLSchemaAttribute(),
					"name": NameSchemaAttribute(),
					"service_id": schema.Int64Attribute{
						Required: true,
					},
					"is_visible": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
				},
			},
		},
		ImportStateCompositeID,
	)
}

type StatusPageMetricWrapper struct {
	upapi.StatusPageMetric

	StatusPageID int64
}

func (w StatusPageMetricWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.PK)
}

type StatusPageMetricResourceModel struct {
	StatusPageID types.Int64  `tfsdk:"statuspage_id"`
	ID           types.Int64  `tfsdk:"id"`
	URL          types.String `tfsdk:"url"`
	Name         types.String `tfsdk:"name"`
	ServiceID    types.Int64  `tfsdk:"service_id"`
	IsVisible    types.Bool   `tfsdk:"is_visible"`
}

func (m StatusPageMetricResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageMetricResourceModelAdapter struct {
	SetAttributeAdapter[int32]
}

func (a StatusPageMetricResourceModelAdapter) Get(
	ctx context.Context, sg StateGetter,
) (*StatusPageMetricResourceModel, diag.Diagnostics) {
	var model StatusPageMetricResourceModel
	if diags := sg.Get(ctx, &model); diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a StatusPageMetricResourceModelAdapter) ToAPIArgument(
	model StatusPageMetricResourceModel,
) (*StatusPageMetricWrapper, error) {
	return &StatusPageMetricWrapper{
		StatusPageID: model.StatusPageID.ValueInt64(),
		StatusPageMetric: upapi.StatusPageMetric{
			PK:        model.ID.ValueInt64(),
			Name:      model.Name.ValueString(),
			ServiceID: model.ServiceID.ValueInt64(),
			IsVisible: model.IsVisible.ValueBool(),
		},
	}, nil
}

func (a StatusPageMetricResourceModelAdapter) FromAPIResult(
	api StatusPageMetricWrapper,
) (*StatusPageMetricResourceModel, error) {
	return &StatusPageMetricResourceModel{
		StatusPageID: types.Int64Value(api.StatusPageID),
		ID:           types.Int64Value(api.PK),
		URL:          types.StringValue(api.URL),
		Name:         types.StringValue(api.Name),
		ServiceID:    types.Int64Value(api.ServiceID),
		IsVisible:    types.BoolValue(api.IsVisible),
	}, nil
}

type StatusPageMetricResourceAPI struct {
	provider *providerImpl
}

func (a StatusPageMetricResourceAPI) Create(ctx context.Context, arg StatusPageMetricWrapper) (*StatusPageMetricWrapper, error) {
	obj, err := a.provider.api.StatusPages().Metrics(upapi.PrimaryKey(arg.StatusPageID)).Create(ctx, arg.StatusPageMetric)
	if err != nil {
		return nil, err
	}

	return &StatusPageMetricWrapper{StatusPageMetric: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageMetricResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*StatusPageMetricWrapper, error) {
	model, ok := arg.(StatusPageMetricResourceModel)
	if !ok {
		return nil, fmt.Errorf("resource read failed:unexpected type %T", arg)
	}

	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	obj, err := a.provider.api.StatusPages().Metrics(statusPageID).Get(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &StatusPageMetricWrapper{StatusPageMetric: *obj, StatusPageID: int64(statusPageID)}, nil
}

func (a StatusPageMetricResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg StatusPageMetricWrapper) (*StatusPageMetricWrapper, error) {
	obj, err := a.provider.api.StatusPages().Metrics(upapi.PrimaryKey(arg.StatusPageID)).Update(ctx, pk, arg.StatusPageMetric)
	if err != nil {
		return nil, err
	}
	return &StatusPageMetricWrapper{StatusPageMetric: *obj, StatusPageID: arg.StatusPageID}, nil
}

func (a StatusPageMetricResourceAPI) Delete(ctx context.Context, arg upapi.PrimaryKeyable) error {
	model, ok := arg.(StatusPageMetricResourceModel)
	if !ok {
		return fmt.Errorf("resource delete failed: unexpected type %T", arg)
	}
	statusPageID := upapi.PrimaryKey(model.StatusPageID.ValueInt64())
	return a.provider.api.StatusPages().Metrics(statusPageID).Delete(ctx, arg)
}
