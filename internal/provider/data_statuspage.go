package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageDataSource{p: p}
}

var StatusPageDataSchema = schema.Schema{
	Description: "Retrieve a list of all status pages configured in your Uptime.com account.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspages": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all status pages in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the status page",
					},
					"url": schema.StringAttribute{
						Computed:    true,
						Description: "API URL for the status page",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the status page",
					},
					"visibility_level": schema.StringAttribute{
						Computed:    true,
						Description: "Visibility level: PUBLIC, UPTIME_USERS, or EXTERNAL_USERS",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Description of the status page",
					},
					"page_type": schema.StringAttribute{
						Computed:    true,
						Description: "Page type: INTERNAL, PUBLIC, or PUBLIC_SLA",
					},
					"slug": schema.StringAttribute{
						Computed:    true,
						Description: "URL slug for the status page",
					},
					"cname": schema.StringAttribute{
						Computed:    true,
						Description: "Custom domain (CNAME) for the status page",
					},
					"timezone": schema.StringAttribute{
						Computed:    true,
						Description: "Timezone for the status page",
					},
					"theme": schema.StringAttribute{
						Computed:    true,
						Description: "Theme for the status page",
					},
				},
			},
		},
	},
}

type StatusPageDataSourceModel struct {
	ID          types.String                    `tfsdk:"id"`
	StatusPages []StatusPageDataSourceItemModel `tfsdk:"statuspages"`
}

type StatusPageDataSourceItemModel struct {
	ID              types.Int64  `tfsdk:"id"`
	URL             types.String `tfsdk:"url"`
	Name            types.String `tfsdk:"name"`
	VisibilityLevel types.String `tfsdk:"visibility_level"`
	Description     types.String `tfsdk:"description"`
	PageType        types.String `tfsdk:"page_type"`
	Slug            types.String `tfsdk:"slug"`
	CNAME           types.String `tfsdk:"cname"`
	Timezone        types.String `tfsdk:"timezone"`
	Theme           types.String `tfsdk:"theme"`
}

var _ datasource.DataSource = &StatusPageDataSource{}

type StatusPageDataSource struct {
	p *providerImpl
}

func (d StatusPageDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspages"
}

func (d StatusPageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageDataSchema
}

func (d StatusPageDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.StatusPages().List(ctx, upapi.StatusPageListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageDataSourceModel{
		ID:          types.StringValue(""),
		StatusPages: make([]StatusPageDataSourceItemModel, len(api)),
	}

	for i := range api {
		model.StatusPages[i] = StatusPageDataSourceItemModel{
			ID:              types.Int64Value(api[i].PK),
			URL:             types.StringValue(api[i].URL),
			Name:            types.StringValue(api[i].Name),
			VisibilityLevel: types.StringValue(api[i].VisibilityLevel),
			Description:     types.StringValue(api[i].Description),
			PageType:        types.StringValue(api[i].PageType),
			Slug:            types.StringValue(api[i].Slug),
			CNAME:           types.StringValue(api[i].CNAME),
			Timezone:        types.StringValue(api[i].Timezone),
			Theme:           types.StringValue(api[i].Theme),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
