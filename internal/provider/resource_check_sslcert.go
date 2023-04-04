package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckSSLCertResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkSSLCertResourceModel, upapi.CheckSSLCert, upapi.Check]{
		api: &checkSSLCertResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_sslcert",
			Schema:         checkSSLCertResourceSchema,
		},
	}
}

var checkSSLCertResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Optional: true,
		},
		"contact_groups": schema.SetAttribute{
			ElementType: types.StringType,
			Required:    true,
		},
		"locations": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"is_paused": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"protocol": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"address": schema.StringAttribute{
			Required: true,
		},
		"port": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"threshold": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"num_retries": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
		"notes": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	},
	Blocks: map[string]schema.Block{
		"config": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"protocol": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"crl": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				"first_element_only": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				"match": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"issuer": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"min_version": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"fingerprint": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"self_signed": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				"url": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}

type checkSSLCertResourceModel struct {
	ID            types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL           types.String `tfsdk:"url" ref:"URL,opt"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	Locations     types.Set    `tfsdk:"locations"`
	Tags          types.Set    `tfsdk:"tags"`
	IsPaused      types.Bool   `tfsdk:"is_paused"`
	Protocol      types.String `tfsdk:"protocol"`
	Address       types.String `tfsdk:"address"`
	Port          types.Int64  `tfsdk:"port"`
	Threshold     types.Int64  `tfsdk:"threshold"`
	NumRetries    types.Int64  `tfsdk:"num_retries"`
	Notes         types.String `tfsdk:"notes"`
	Config        types.Object `tfsdk:"config" ref:"SSLConfig"`
}

var _ genericResourceAPI[upapi.CheckSSLCert, upapi.Check] = (*checkSSLCertResourceAPI)(nil)

type checkSSLCertResourceAPI struct {
	provider *providerImpl
}

func (c *checkSSLCertResourceAPI) Create(ctx context.Context, arg upapi.CheckSSLCert) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateSSLCert(ctx, arg)
}

func (c *checkSSLCertResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c *checkSSLCertResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckSSLCert) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateSSLCert(ctx, pk, arg)
}

func (c *checkSSLCertResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
