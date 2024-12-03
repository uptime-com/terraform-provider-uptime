package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[StatusPageResourceModel, upapi.StatusPage, upapi.StatusPage]{
		api: &StatusPageResourceAPI{provider: p},
		mod: StatusPageResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "statuspage",
			Schema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":   IDSchemaAttribute(),
					"url":  URLSchemaAttribute(),
					"name": NameSchemaAttribute(),
					"visibility_level": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("UPTIME_USERS"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"PUBLIC", "UPTIME_USERS", "EXTERNAL_USERS"}),
						},
					},
					"description": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"page_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("INTERNAL"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"INTERNAL", "PUBLIC", "PUBLIC_SLA"}),
						},
					},
					"slug": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"cname": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
						Validators: []validator.String{
							HostnameValidator(),
						},
					},
					"allow_subscriptions": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"allow_search_indexing": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"allow_drill_down": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"auth_username": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"auth_password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
						Computed:  true,
						Default:   stringdefault.StaticString(""),
					},
					"show_status_tab": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"show_active_incidents": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"show_component_response_time": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"show_history_tab": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"default_history_date_range": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(90),
					},
					"uptime_calculation_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("BY_INCIDENTS"),
					},
					"show_history_snake": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"show_component_history": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"show_summary_metrics": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"show_past_incidents": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"allow_pdf_report": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"google_analytics_code": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"contact_email": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"email_from": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"email_reply_to": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"custom_header_html": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"custom_footer_html": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"custom_css": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"company_website_url": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"timezone": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("GMT"),
					},
					"allow_subscriptions_email": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"allow_subscriptions_rss": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"allow_subscriptions_slack": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"allow_subscriptions_sms": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"allow_subscriptions_webhook": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"hide_empty_tabs_history": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"theme": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("LEGACY"),
					},
					"custom_header_bg_color_hex": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#002E52"),
					},
					"custom_header_text_color_hex": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("#FFFFFF"),
					},
				},
			},
		},
	}
}

type StatusPageResourceModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	URL                       types.String `tfsdk:"url"`
	Name                      types.String `tfsdk:"name"`
	VisibilityLevel           types.String `tfsdk:"visibility_level"`
	Description               types.String `tfsdk:"description"`
	PageType                  types.String `tfsdk:"page_type"`
	Slug                      types.String `tfsdk:"slug"`
	CNAME                     types.String `tfsdk:"cname"`
	AllowSubscriptions        types.Bool   `tfsdk:"allow_subscriptions"`
	AllowSearchIndexing       types.Bool   `tfsdk:"allow_search_indexing"`
	AllowDrillDown            types.Bool   `tfsdk:"allow_drill_down"`
	AuthUsername              types.String `tfsdk:"auth_username"`
	AuthPassword              types.String `tfsdk:"auth_password"`
	ShowStatusTab             types.Bool   `tfsdk:"show_status_tab"`
	ShowActiveIncidents       types.Bool   `tfsdk:"show_active_incidents"`
	ShowComponentResponseTime types.Bool   `tfsdk:"show_component_response_time"`
	ShowHistoryTab            types.Bool   `tfsdk:"show_history_tab"`
	DefaultHistoryDateRange   types.Int64  `tfsdk:"default_history_date_range"`
	UptimeCalculationType     types.String `tfsdk:"uptime_calculation_type"`
	ShowHistorySnake          types.Bool   `tfsdk:"show_history_snake"`
	ShowComponentHistory      types.Bool   `tfsdk:"show_component_history"`
	ShowSummaryMetrics        types.Bool   `tfsdk:"show_summary_metrics"`
	ShowPastIncidents         types.Bool   `tfsdk:"show_past_incidents"`
	AllowPdfReport            types.Bool   `tfsdk:"allow_pdf_report"`
	GoogleAnalyticsCode       types.String `tfsdk:"google_analytics_code"`
	ContactEmail              types.String `tfsdk:"contact_email"`
	EmailFrom                 types.String `tfsdk:"email_from"`
	EmailReplyTo              types.String `tfsdk:"email_reply_to"`
	CustomHeaderHtml          types.String `tfsdk:"custom_header_html"`
	CustomFooterHtml          types.String `tfsdk:"custom_footer_html"`
	CustomCss                 types.String `tfsdk:"custom_css"`
	CompanyWebsiteUrl         types.String `tfsdk:"company_website_url"`
	Timezone                  types.String `tfsdk:"timezone"`
	AllowSubscriptionsEmail   types.Bool   `tfsdk:"allow_subscriptions_email"`
	AllowSubscriptionsRss     types.Bool   `tfsdk:"allow_subscriptions_rss"`
	AllowSubscriptionsSlack   types.Bool   `tfsdk:"allow_subscriptions_slack"`
	AllowSubscriptionsSms     types.Bool   `tfsdk:"allow_subscriptions_sms"`
	AllowSubscriptionsWebhook types.Bool   `tfsdk:"allow_subscriptions_webhook"`
	HideEmptyTabsHistory      types.Bool   `tfsdk:"hide_empty_tabs_history"`
	Theme                     types.String `tfsdk:"theme"`
	CustomHeaderBgColorHex    types.String `tfsdk:"custom_header_bg_color_hex"`
	CustomHeaderTextColorHex  types.String `tfsdk:"custom_header_text_color_hex"`
}

func (m StatusPageResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type StatusPageResourceModelAdapter struct{}

func (c StatusPageResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*StatusPageResourceModel, diag.Diagnostics) {
	model := *new(StatusPageResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (c StatusPageResourceModelAdapter) ToAPIArgument(model StatusPageResourceModel) (*upapi.StatusPage, error) {
	api := upapi.StatusPage{
		Name:                      model.Name.ValueString(),
		VisibilityLevel:           model.VisibilityLevel.ValueString(),
		Description:               model.Description.ValueString(),
		PageType:                  model.PageType.ValueString(),
		Slug:                      model.Slug.ValueString(),
		CNAME:                     model.CNAME.ValueString(),
		AllowSubscriptions:        model.AllowSubscriptions.ValueBool(),
		AllowSearchIndexing:       model.AllowSearchIndexing.ValueBool(),
		AllowDrillDown:            model.AllowDrillDown.ValueBool(),
		AuthUsername:              model.AuthUsername.ValueString(),
		AuthPassword:              model.AuthPassword.ValueString(),
		ShowStatusTab:             model.ShowStatusTab.ValueBool(),
		ShowActiveIncidents:       model.ShowActiveIncidents.ValueBool(),
		ShowComponentResponseTime: model.ShowComponentResponseTime.ValueBool(),
		ShowHistoryTab:            model.ShowHistoryTab.ValueBool(),
		DefaultHistoryDateRange:   model.DefaultHistoryDateRange.ValueInt64(),
		UptimeCalculationType:     model.UptimeCalculationType.ValueString(),
		ShowHistorySnake:          model.ShowHistorySnake.ValueBool(),
		ShowComponentHistory:      model.ShowComponentHistory.ValueBool(),
		ShowSummaryMetrics:        model.ShowSummaryMetrics.ValueBool(),
		ShowPastIncidents:         model.ShowPastIncidents.ValueBool(),
		AllowPdfReport:            model.AllowPdfReport.ValueBool(),
		GoogleAnalyticsCode:       model.GoogleAnalyticsCode.ValueString(),
		ContactEmail:              model.ContactEmail.ValueString(),
		EmailFrom:                 model.EmailFrom.ValueString(),
		EmailReplyTo:              model.EmailReplyTo.ValueString(),
		CustomHeaderHtml:          model.CustomHeaderHtml.ValueString(),
		CustomFooterHtml:          model.CustomFooterHtml.ValueString(),
		CustomCss:                 model.CustomCss.ValueString(),
		CompanyWebsiteUrl:         model.CompanyWebsiteUrl.ValueString(),
		Timezone:                  model.Timezone.ValueString(),
		AllowSubscriptionsEmail:   model.AllowSubscriptionsEmail.ValueBool(),
		AllowSubscriptionsRss:     model.AllowSubscriptionsRss.ValueBool(),
		AllowSubscriptionsSlack:   model.AllowSubscriptionsSlack.ValueBool(),
		AllowSubscriptionsSms:     model.AllowSubscriptionsSms.ValueBool(),
		AllowSubscriptionsWebhook: model.AllowSubscriptionsWebhook.ValueBool(),
		HideEmptyTabsHistory:      model.HideEmptyTabsHistory.ValueBool(),
		Theme:                     model.Theme.ValueString(),
		CustomHeaderBgColorHex:    model.CustomHeaderBgColorHex.ValueString(),
		CustomHeaderTextColorHex:  model.CustomHeaderTextColorHex.ValueString(),
	}
	return &api, nil
}

func (c StatusPageResourceModelAdapter) FromAPIResult(api upapi.StatusPage) (*StatusPageResourceModel, error) {
	model := StatusPageResourceModel{
		ID:                        types.Int64Value(api.PK),
		URL:                       types.StringValue(api.URL),
		Name:                      types.StringValue(api.Name),
		VisibilityLevel:           types.StringValue(api.VisibilityLevel),
		Description:               types.StringValue(api.Description),
		PageType:                  types.StringValue(api.PageType),
		Slug:                      types.StringValue(api.Slug),
		CNAME:                     types.StringValue(api.CNAME),
		AllowSubscriptions:        types.BoolValue(api.AllowSubscriptions),
		AllowSearchIndexing:       types.BoolValue(api.AllowSearchIndexing),
		AllowDrillDown:            types.BoolValue(api.AllowDrillDown),
		AuthUsername:              types.StringValue(api.AuthUsername),
		AuthPassword:              types.StringValue(api.AuthPassword),
		ShowStatusTab:             types.BoolValue(api.ShowStatusTab),
		ShowActiveIncidents:       types.BoolValue(api.ShowActiveIncidents),
		ShowComponentResponseTime: types.BoolValue(api.ShowComponentResponseTime),
		ShowHistoryTab:            types.BoolValue(api.ShowHistoryTab),
		DefaultHistoryDateRange:   types.Int64Value(api.DefaultHistoryDateRange),
		UptimeCalculationType:     types.StringValue(api.UptimeCalculationType),
		ShowHistorySnake:          types.BoolValue(api.ShowHistorySnake),
		ShowComponentHistory:      types.BoolValue(api.ShowComponentHistory),
		ShowSummaryMetrics:        types.BoolValue(api.ShowSummaryMetrics),
		ShowPastIncidents:         types.BoolValue(api.ShowPastIncidents),
		AllowPdfReport:            types.BoolValue(api.AllowPdfReport),
		GoogleAnalyticsCode:       types.StringValue(api.GoogleAnalyticsCode),
		ContactEmail:              types.StringValue(api.ContactEmail),
		EmailFrom:                 types.StringValue(api.EmailFrom),
		EmailReplyTo:              types.StringValue(api.EmailReplyTo),
		CustomHeaderHtml:          types.StringValue(api.CustomHeaderHtml),
		CustomFooterHtml:          types.StringValue(api.CustomFooterHtml),
		CustomCss:                 types.StringValue(api.CustomCss),
		CompanyWebsiteUrl:         types.StringValue(api.CompanyWebsiteUrl),
		Timezone:                  types.StringValue(api.Timezone),
		AllowSubscriptionsEmail:   types.BoolValue(api.AllowSubscriptionsEmail),
		AllowSubscriptionsRss:     types.BoolValue(api.AllowSubscriptionsRss),
		AllowSubscriptionsSlack:   types.BoolValue(api.AllowSubscriptionsSlack),
		AllowSubscriptionsSms:     types.BoolValue(api.AllowSubscriptionsSms),
		AllowSubscriptionsWebhook: types.BoolValue(api.AllowSubscriptionsWebhook),
		HideEmptyTabsHistory:      types.BoolValue(api.HideEmptyTabsHistory),
		Theme:                     types.StringValue(api.Theme),
		CustomHeaderBgColorHex:    types.StringValue(api.CustomHeaderBgColorHex),
		CustomHeaderTextColorHex:  types.StringValue(api.CustomHeaderTextColorHex),
	}
	return &model, nil
}

type StatusPageResourceAPI struct {
	provider *providerImpl
}

func (s StatusPageResourceAPI) Create(ctx context.Context, arg upapi.StatusPage) (*upapi.StatusPage, error) {
	obj, err := s.provider.api.StatusPages().Create(ctx, arg)
	return obj, err
}

func (s StatusPageResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*upapi.StatusPage, error) {
	obj, err := s.provider.api.StatusPages().Get(ctx, arg)
	return obj, err
}

func (s StatusPageResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.StatusPage) (*upapi.StatusPage, error) {
	obj, err := s.provider.api.StatusPages().Update(ctx, pk, arg)
	return obj, err
}

func (s StatusPageResourceAPI) Delete(ctx context.Context, keyable upapi.PrimaryKeyable) error {
	return s.provider.api.StatusPages().Delete(ctx, keyable)
}
