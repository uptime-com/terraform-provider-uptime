package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		"id":             IDAttribute(),
		"url":            URLAttribute(),
		"name":           NameAttribute(),
		"contact_groups": ContactGroupsAttribute(),
		"locations":      LocationsReadOnlyAttribute(),
		"tags":           TagsAttribute(),
		"is_paused":      IsPausedAttribute(),
		"threshold": ThresholdDescriptionAttribute(20,
			"Raise an alert if there are less than this many days before the SSL certificate needs to be renewed"),
		"num_retries": NumRetriesAttribute(2),
		"notes":       NotesAttribute(),

		"address": schema.StringAttribute{
			Required: true,
		},
		"port": schema.Int64Attribute{
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(443),
		},

		"config": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"protocol": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("https"),
					Description: "Application level protocol",
					Validators: []validator.String{
						OneOfStringValidator([]string{
							"https", "ftp", "ftps", "http", "h2", "imap", "imaps", "irc", "ircs", "ldap", "ldaps", "mysql",
							"pop3", "pop3s", "postgres", "sieve", "smtp", "smtps", "xmpp", "xmpp-server",
						}),
					},
				},
				"crl": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
				"first_element_only": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
				"match": schema.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString(""),
				},
				"issuer": schema.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString(""),
				},
				"min_version": schema.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString("sslv3"),
					Validators: []validator.String{
						OneOfStringValidator([]string{"sslv3", "tlsv1", "tlsv11", "tlsv12", "tlsv13"}),
					},
				},
				"fingerprint": schema.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  stringdefault.StaticString(""),
				},
				"self_signed": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
				"url": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "Specify location of certificate or CRL file by URL, instead of retrieving from main domain address.",
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
