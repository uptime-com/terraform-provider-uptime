package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewScheduledReportsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return ScheduledReportsDataSource{p: p}
}

// ScheduledReportsDataSchema defines the schema for the scheduled reports data source.
var ScheduledReportsDataSchema = schema.Schema{
	Description: "Retrieve a list of all scheduled reports configured in your Uptime.com account. Scheduled reports automatically send SLA reports to recipients.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"scheduled_reports": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all scheduled reports in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the scheduled report",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the scheduled report resource",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the scheduled report",
					},
					"sla_report": schema.StringAttribute{
						Computed:    true,
						Description: "URL of the associated SLA report",
					},
					"recipient_users": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of user URLs to receive the report",
					},
					"recipient_emails": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of email addresses to receive the report",
					},
					"file_type": schema.StringAttribute{
						Computed:    true,
						Description: "File type for the report (e.g., 'pdf', 'csv')",
					},
					"recurrence": schema.StringAttribute{
						Computed:    true,
						Description: "Recurrence pattern (e.g., 'daily', 'weekly', 'monthly')",
					},
					"on_weekday": schema.Int64Attribute{
						Computed:    true,
						Description: "Day of week for weekly reports (0=Monday, 6=Sunday)",
					},
					"at_time": schema.Int64Attribute{
						Computed:    true,
						Description: "Time of day to send report (in minutes from midnight)",
					},
					"is_enabled": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the scheduled report is enabled",
					},
				},
			},
		},
	},
}

type ScheduledReportsDataSourceModel struct {
	ID               types.String                          `tfsdk:"id"`
	ScheduledReports []ScheduledReportsDataSourceItemModel `tfsdk:"scheduled_reports"`
}

type ScheduledReportsDataSourceItemModel struct {
	ID              types.Int64  `tfsdk:"id"`
	URL             types.String `tfsdk:"url"`
	Name            types.String `tfsdk:"name"`
	SLAReport       types.String `tfsdk:"sla_report"`
	RecipientUsers  types.List   `tfsdk:"recipient_users"`
	RecipientEmails types.List   `tfsdk:"recipient_emails"`
	FileType        types.String `tfsdk:"file_type"`
	Recurrence      types.String `tfsdk:"recurrence"`
	OnWeekday       types.Int64  `tfsdk:"on_weekday"`
	AtTime          types.Int64  `tfsdk:"at_time"`
	IsEnabled       types.Bool   `tfsdk:"is_enabled"`
}

var _ datasource.DataSource = &ScheduledReportsDataSource{}

type ScheduledReportsDataSource struct {
	p *providerImpl
}

func (d ScheduledReportsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_scheduled_reports"
}

func (d ScheduledReportsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = ScheduledReportsDataSchema
}

func (d ScheduledReportsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.ScheduledReports().List(ctx, upapi.ScheduledReportListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := ScheduledReportsDataSourceModel{
		ID:               types.StringValue(""),
		ScheduledReports: make([]ScheduledReportsDataSourceItemModel, len(api.Items)),
	}

	for i := range api.Items {
		// Convert RecipientUsers slice to types.List
		recipientUsersValues := make([]attr.Value, len(api.Items[i].RecipientUsers))
		for j, v := range api.Items[i].RecipientUsers {
			recipientUsersValues[j] = types.StringValue(v)
		}
		recipientUsers := types.ListNull(types.StringType)
		if len(recipientUsersValues) > 0 {
			recipientUsers = types.ListValueMust(types.StringType, recipientUsersValues)
		}

		// Convert RecipientEmails slice to types.List
		recipientEmailsValues := make([]attr.Value, len(api.Items[i].RecipientEmails))
		for j, v := range api.Items[i].RecipientEmails {
			recipientEmailsValues[j] = types.StringValue(v)
		}
		recipientEmails := types.ListNull(types.StringType)
		if len(recipientEmailsValues) > 0 {
			recipientEmails = types.ListValueMust(types.StringType, recipientEmailsValues)
		}

		model.ScheduledReports[i] = ScheduledReportsDataSourceItemModel{
			ID:              types.Int64Value(api.Items[i].PK),
			URL:             types.StringValue(api.Items[i].URL),
			Name:            types.StringValue(api.Items[i].Name),
			SLAReport:       types.StringValue(api.Items[i].ScheduledReport),
			RecipientUsers:  recipientUsers,
			RecipientEmails: recipientEmails,
			FileType:        types.StringValue(api.Items[i].FileType),
			Recurrence:      types.StringValue(api.Items[i].Recurrence),
			OnWeekday:       types.Int64Value(int64(api.Items[i].OnWeekday)),
			AtTime:          types.Int64Value(int64(api.Items[i].AtTime)),
			IsEnabled:       types.BoolValue(api.Items[i].IsEnabled),
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
