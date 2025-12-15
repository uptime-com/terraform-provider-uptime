package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewStatusPageSubscriberDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return StatusPageSubscriberDataSource{p: p}
}

var StatusPageSubscriberDataSchema = schema.Schema{
	Description: "Retrieve a list of all subscribers for a specific status page.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"statuspage_id": schema.Int64Attribute{
			Required:    true,
			Description: "ID of the status page to retrieve subscribers for",
		},
		"subscribers": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all subscribers for the status page",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the subscriber",
					},
					"target": schema.StringAttribute{
						Computed:    true,
						Description: "Target (email, phone, etc.) for notifications",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Type of subscription (email, sms, webhook, etc.)",
					},
					"force_validation_sms": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether SMS validation is forced",
					},
				},
			},
		},
	},
}

type StatusPageSubscriberDataSourceModel struct {
	ID           types.String                              `tfsdk:"id"`
	StatusPageID types.Int64                               `tfsdk:"statuspage_id"`
	Subscribers  []StatusPageSubscriberDataSourceItemModel `tfsdk:"subscribers"`
}

type StatusPageSubscriberDataSourceItemModel struct {
	ID                 types.Int64  `tfsdk:"id"`
	Target             types.String `tfsdk:"target"`
	Type               types.String `tfsdk:"type"`
	ForceValidationSMS types.Bool   `tfsdk:"force_validation_sms"`
}

var _ datasource.DataSource = &StatusPageSubscriberDataSource{}

type StatusPageSubscriberDataSource struct {
	p *providerImpl
}

func (d StatusPageSubscriberDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_statuspage_subscribers"
}

func (d StatusPageSubscriberDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = StatusPageSubscriberDataSchema
}

func (d StatusPageSubscriberDataSource) Read(ctx context.Context, rq datasource.ReadRequest, rs *datasource.ReadResponse) {
	var config StatusPageSubscriberDataSourceModel
	diags := rq.Config.Get(ctx, &config)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	pk := upapi.PrimaryKey(config.StatusPageID.ValueInt64())
	api, err := d.p.api.StatusPages().Subscribers(pk).List(ctx, upapi.StatusPageSubscriberListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := StatusPageSubscriberDataSourceModel{
		ID:           types.StringValue(""),
		StatusPageID: config.StatusPageID,
		Subscribers:  make([]StatusPageSubscriberDataSourceItemModel, len(api.Items)),
	}

	for i := range api.Items {
		model.Subscribers[i] = StatusPageSubscriberDataSourceItemModel{
			ID:                 types.Int64Value(api.Items[i].PK),
			Target:             types.StringValue(api.Items[i].Target),
			Type:               types.StringValue(api.Items[i].Type),
			ForceValidationSMS: types.BoolValue(api.Items[i].ForceValidationSMS),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
