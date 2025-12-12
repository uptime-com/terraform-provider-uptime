package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shopspring/decimal"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckWebhookResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[CheckWebhookResourceModel, upapi.CheckWebhook, upapi.Check](
		CheckWebookResourceAPI{provider: p},
		CheckWebhookResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_webhook",
			Schema: schema.Schema{
				Description: "Receive alerts based on periodic jobs or processes using an automated HTTP callback. Import using the check ID: `terraform import uptime_check_webhook.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":                        IDSchemaAttribute(),
					"url":                       URLSchemaAttribute(),
					"name":                      NameSchemaAttribute(),
					"contact_groups":            ContactGroupsSchemaAttribute(),
					"tags":                      TagsSchemaAttribute(),
					"is_paused":                 IsPausedSchemaAttribute(),
					"notes":                     NotesSchemaAttribute(),
					"sla":                       SLASchemaAttribute(),
					"include_in_global_metrics": IncludeInGlobalMetricsSchemaAttribute(),
					"webhook_url": schema.StringAttribute{
						Computed:    true,
						Description: "URL to send data to your check",
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type CheckWebhookResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url"  ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Notes                  types.String `tfsdk:"notes"`
	SLA                    types.Object `tfsdk:"sla"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	WebhookURL             types.String `tfsdk:"webhook_url"`

	sla *SLAAttribute `tfsdk:"-"`
}

func (m CheckWebhookResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckWebhookResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckWebhookResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckWebhookResourceModel, diag.Diagnostics) {
	model := *new(CheckWebhookResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.sla, diags = a.SLAAttributeContext(ctx, model.SLA)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckWebhookResourceModelAdapter) ToAPIArgument(model CheckWebhookResourceModel) (*upapi.CheckWebhook, error) {
	api := upapi.CheckWebhook{
		Name:                   model.Name.ValueString(),
		ContactGroups:          a.ContactGroups(model.ContactGroups),
		Tags:                   a.Tags(model.Tags),
		IsPaused:               model.IsPaused.ValueBool(),
		Notes:                  model.Notes.ValueString(),
		IncludeInGlobalMetrics: model.IncludeInGlobalMetrics.ValueBool(),
		WebhookUrl:             model.WebhookURL.ValueString(),
	}
	if model.sla != nil {
		if !model.sla.Uptime.IsUnknown() {
			api.UptimeSLA = model.sla.Uptime.ValueDecimal()
		}
		if !model.sla.Latency.IsUnknown() {
			api.ResponseTimeSLA = decimal.NewFromFloat(model.sla.Latency.ValueDuration().Seconds())
		}
	}

	return &api, nil
}

func (a CheckWebhookResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckWebhookResourceModel, error) {
	model := CheckWebhookResourceModel{
		ID:                     types.Int64Value(api.PK),
		URL:                    types.StringValue(api.URL),
		Name:                   types.StringValue(api.Name),
		ContactGroups:          a.ContactGroupsValue(api.ContactGroups),
		Tags:                   a.TagsValue(api.Tags),
		IsPaused:               types.BoolValue(api.IsPaused),
		Notes:                  types.StringValue(api.Notes),
		IncludeInGlobalMetrics: types.BoolValue(api.IncludeInGlobalMetrics),
		WebhookURL:             types.StringValue(api.WebhookURL),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Latency: DurationValueFromDecimalSeconds(api.ResponseTimeSLA),
			Uptime:  DecimalValue(api.UptimeSLA),
		}),
	}
	return &model, nil
}

type CheckWebookResourceAPI struct {
	provider *providerImpl
}

func (c CheckWebookResourceAPI) Create(ctx context.Context, arg upapi.CheckWebhook) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateWebhook(ctx, arg)
}

func (c CheckWebookResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckWebookResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckWebhook) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateWebhook(ctx, pk, arg)
}

func (c CheckWebookResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
