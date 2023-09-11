package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[statusPageResourceModel, upapi.StatusPage, upapi.StatusPage]{
		api: &statusPageResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "statuspage",
			Schema:         statusPageResourceSchema,
		},
	}
}

type statusPageResourceModel struct {
	ID                        types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                       types.String `tfsdk:"url"  ref:"URL,opt"`
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
}

var statusPageResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"visibility_level": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("UPTIME_USERS"),
		},
		"description": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(""),
		},
		"page_type": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("PUBLIC"),
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
			Default:  stringdefault.StaticString("UTC"),
		},
	},
}

var _ genericResourceAPI[upapi.StatusPage, upapi.StatusPage] = (*statusPageResourceAPI)(nil)

type statusPageResourceAPI struct {
	provider *providerImpl
}

func (s *statusPageResourceAPI) Create(ctx context.Context, arg upapi.StatusPage) (*upapi.StatusPage, error) {
	obj, err := s.provider.api.StatusPages().Create(ctx, arg)
	return obj, err
}

func (s *statusPageResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*upapi.StatusPage, error) {
	obj, err := s.provider.api.StatusPages().Get(ctx, arg)
	return obj, err
}

func (s *statusPageResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.StatusPage) (*upapi.StatusPage, error) {
	obj, err := s.provider.api.StatusPages().Update(ctx, pk, arg)
	return obj, err
}

func (s *statusPageResourceAPI) Delete(ctx context.Context, keyable upapi.PrimaryKeyable) error {
	return s.provider.api.StatusPages().Delete(ctx, keyable)
}
