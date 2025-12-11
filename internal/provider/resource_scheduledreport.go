package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewScheduledReportResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[ScheduledReportResourceModel, upapi.ScheduledReport, upapi.ScheduledReport](
		&ScheduledReportResourceAPI{provider: p},
		ScheduledReportResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "scheduled_report",
			Schema: schema.Schema{
				Description: "Scheduled report resource. Import using the scheduled report ID: `terraform import uptime_scheduled_report.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":   IDSchemaAttribute(),
					"url":  URLSchemaAttribute(),
					"name": NameSchemaAttribute(),
					"sla_report": schema.StringAttribute{
						Required:    true,
						Description: "Select an SLA report to send on this schedule",
					},
					"recipient_users": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
						Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"recipient_emails": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
						Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"file_type": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("PDF"),
						Description: "Report file type, valid values are PDF(default) or XLS",
						Validators: []validator.String{
							OneOfStringValidator([]string{"PDF", "XLS"}),
						},
					},
					"recurrence": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("DAILY"),
						Description: "How often to deliver this report. Valid values are DAILY, WEEKLY, MONTHLY, QUARTERLY, YEARLY",
						Validators: []validator.String{
							OneOfStringValidator([]string{"DAILY", "WEEKLY", "MONTHLY", "QUARTERLY", "YEARLY"}),
						},
					},
					"on_weekday": schema.Int32Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int32default.StaticInt32(1),
						Description: "Weekly reports will be sent on this day",
						Validators: []validator.Int32{
							int32validator.Between(1, 7),
						},
					},
					"at_time": schema.Int32Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int32default.StaticInt32(1),
						Description: "Reports will be sent at this time (local time), value is hour of day: 0-23",
						Validators: []validator.Int32{
							int32validator.Between(0, 23),
						},
					},
					"is_enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
				},
			},
		},
		ImportStateSimpleID,
	)
}

type ScheduledReportResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	URL             types.String `tfsdk:"url"`
	Name            types.String `tfsdk:"name"`
	ScheduledReport types.String `tfsdk:"sla_report"`
	RecipientUsers  types.Set    `tfsdk:"recipient_users"`
	RecipientEmails types.Set    `tfsdk:"recipient_emails"`
	FileType        types.String `tfsdk:"file_type"`
	Recurrence      types.String `tfsdk:"recurrence"`
	OnWeekday       types.Int32  `tfsdk:"on_weekday"`
	AtTime          types.Int32  `tfsdk:"at_time"`
	IsEnabled       types.Bool   `tfsdk:"is_enabled"`
}

func (m ScheduledReportResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type ScheduledReportResourceModelAdapter struct {
	TagsAttributeAdapter
	SetAttributeAdapter[string]
}

func (a ScheduledReportResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*ScheduledReportResourceModel, diag.Diagnostics) {
	model := *new(ScheduledReportResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a ScheduledReportResourceModelAdapter) ToAPIArgument(model ScheduledReportResourceModel) (*upapi.ScheduledReport, error) {
	return &upapi.ScheduledReport{
		Name:            model.Name.ValueString(),
		ScheduledReport: model.ScheduledReport.ValueString(),
		RecipientUsers:  a.Slice(model.RecipientUsers),
		RecipientEmails: a.Slice(model.RecipientEmails),
		FileType:        model.FileType.ValueString(),
		Recurrence:      model.Recurrence.ValueString(),
		OnWeekday:       model.OnWeekday.ValueInt32(),
		AtTime:          model.AtTime.ValueInt32(),
		IsEnabled:       model.IsEnabled.ValueBool(),
	}, nil
}

func (a ScheduledReportResourceModelAdapter) FromAPIResult(api upapi.ScheduledReport) (*ScheduledReportResourceModel, error) {
	return &ScheduledReportResourceModel{
		ID:              types.Int64Value(api.PK),
		URL:             types.StringValue(api.URL),
		Name:            types.StringValue(api.Name),
		ScheduledReport: types.StringValue(api.ScheduledReport),
		RecipientUsers:  a.SliceValue(api.RecipientUsers),
		RecipientEmails: a.SliceValue(api.RecipientEmails),
		FileType:        types.StringValue(api.FileType),
		Recurrence:      types.StringValue(api.Recurrence),
		OnWeekday:       types.Int32Value(int32(api.OnWeekday)),
		AtTime:          types.Int32Value(int32(api.AtTime)),
		IsEnabled:       types.BoolValue(api.IsEnabled),
	}, nil
}

type ScheduledReportResourceAPI struct {
	provider *providerImpl
}

func (a ScheduledReportResourceAPI) Create(ctx context.Context, arg upapi.ScheduledReport) (*upapi.ScheduledReport, error) {
	obj, err := a.provider.api.ScheduledReports().Create(ctx, arg)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (a ScheduledReportResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*upapi.ScheduledReport, error) {
	return a.provider.api.ScheduledReports().Get(ctx, arg)
}

func (a ScheduledReportResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.ScheduledReport) (*upapi.ScheduledReport, error) {
	return a.provider.api.ScheduledReports().Update(ctx, pk, arg)
}

func (a ScheduledReportResourceAPI) Delete(ctx context.Context, keyable upapi.PrimaryKeyable) error {
	return a.provider.api.ScheduledReports().Delete(ctx, keyable)
}
