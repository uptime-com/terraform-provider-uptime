package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckSSLCertResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckSSLCertResourceModel, upapi.CheckSSLCert, upapi.Check]{
		api: CheckSSLCertResourceAPI{provider: p},
		mod: CheckSSLCertResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_sslcert",
			Schema: schema.Schema{
				Description: "Verify SSL certificate validity",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"locations":      LocationsReadOnlySchemaAttribute(),
					"tags":           TagsSchemaAttribute(),
					"is_paused":      IsPausedSchemaAttribute(),
					"threshold": ThresholdDescriptionSchemaAttribute(
						20,
						"Raise an alert if there are less than this many days before the SSL certificate needs to be renewed",
					),
					"num_retries": NumRetriesSchemaAttribute(2),
					"notes":       NotesSchemaAttribute(),
					"address":     AddressHostnameSchemaAttribute(),
					"port":        PortSchemaAttribute(443),

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
								// This is workaround for API bug where it fails to update the value of this attribute once it has been set to true
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.RequiresReplace(),
								},
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
			},
		},
	}
}

type CheckSSLCertResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
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
	Config        types.Object `tfsdk:"config"`

	config *CheckSSLCertConfigAttribute `tfsdk:"-"`
}

func (m CheckSSLCertResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckSSLCertConfigAttribute struct {
	Protocol         types.String `tfsdk:"protocol"`
	CRL              types.Bool   `tfsdk:"crl"`
	FirstElementOnly types.Bool   `tfsdk:"first_element_only"`
	Match            types.String `tfsdk:"match"`
	Issuer           types.String `tfsdk:"issuer"`
	MinVersion       types.String `tfsdk:"min_version"`
	Fingerprint      types.String `tfsdk:"fingerprint"`
	SelfSigned       types.Bool   `tfsdk:"self_signed"`
	URL              types.String `tfsdk:"url"`
}

type CheckSSLCertResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
}

func (a CheckSSLCertResourceModelAdapter) ConfigAttributeContext(ctx context.Context, v types.Object) (*CheckSSLCertConfigAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := *new(CheckSSLCertConfigAttribute)
	diags := v.As(ctx, &m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (a CheckSSLCertResourceModelAdapter) configAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"protocol":           types.StringType,
		"crl":                types.BoolType,
		"first_element_only": types.BoolType,
		"match":              types.StringType,
		"issuer":             types.StringType,
		"min_version":        types.StringType,
		"fingerprint":        types.StringType,
		"self_signed":        types.BoolType,
		"url":                types.StringType,
	}
}

func (a CheckSSLCertResourceModelAdapter) configAttributeValues(model CheckSSLCertConfigAttribute) map[string]attr.Value {
	return map[string]attr.Value{
		"protocol":           model.Protocol,
		"crl":                model.CRL,
		"first_element_only": model.FirstElementOnly,
		"match":              model.Match,
		"issuer":             model.Issuer,
		"min_version":        model.MinVersion,
		"fingerprint":        model.Fingerprint,
		"self_signed":        model.SelfSigned,
		"url":                model.URL,
	}
}

func (a CheckSSLCertResourceModelAdapter) ConfigAttributeValue(m CheckSSLCertConfigAttribute) types.Object {
	return types.ObjectValueMust(a.configAttributeTypes(), a.configAttributeValues(m))
}

func (a CheckSSLCertResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckSSLCertResourceModel, diag.Diagnostics) {
	model := *new(CheckSSLCertResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.config, diags = a.ConfigAttributeContext(ctx, model.Config)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckSSLCertResourceModelAdapter) ToAPIArgument(model CheckSSLCertResourceModel) (*upapi.CheckSSLCert, error) {
	api := upapi.CheckSSLCert{
		Name:          model.Name.ValueString(),
		Address:       model.Address.ValueString(),
		Port:          model.Port.ValueInt64(),
		ContactGroups: a.ContactGroups(model.ContactGroups),
		Locations:     a.Locations(model.Locations),
		Tags:          a.Tags(model.Tags),
		IsPaused:      model.IsPaused.ValueBool(),
		Threshold:     model.Threshold.ValueInt64(),
		NumRetries:    model.NumRetries.ValueInt64(),
		Notes:         model.Notes.ValueString(),
	}
	if model.config != nil {
		api.SSLConfig = upapi.CheckSSLCertConfig{
			Protocol:         model.config.Protocol.ValueString(),
			CRL:              model.config.CRL.ValueBool(),
			FirstElementOnly: model.config.FirstElementOnly.ValueBool(),
			Match:            model.config.Match.ValueString(),
			Issuer:           model.config.Issuer.ValueString(),
			MinVersion:       model.config.MinVersion.ValueString(),
			Fingerprint:      model.config.Fingerprint.ValueString(),
			SelfSigned:       model.config.SelfSigned.ValueBool(),
			URL:              model.config.URL.ValueString(),
		}
	}
	return &api, nil
}

func (a CheckSSLCertResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckSSLCertResourceModel, error) {
	model := CheckSSLCertResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsValue(api.ContactGroups),
		Locations:     a.LocationsValue(api.Locations),
		Tags:          a.TagsValue(api.Tags),
		IsPaused:      types.BoolValue(api.IsPaused),
		Address:       types.StringValue(api.Address),
		Port:          types.Int64Value(api.Port),
		Threshold:     types.Int64Value(api.Threshold),
		NumRetries:    types.Int64Value(api.NumRetries),
		Notes:         types.StringValue(api.Notes),
		Config: a.ConfigAttributeValue(CheckSSLCertConfigAttribute{
			Protocol:         types.StringValue(api.SSLConfig.Protocol),
			CRL:              types.BoolValue(api.SSLConfig.CRL),
			FirstElementOnly: types.BoolValue(api.SSLConfig.FirstElementOnly),
			Match:            types.StringValue(api.SSLConfig.Match),
			Issuer:           types.StringValue(api.SSLConfig.Issuer),
			MinVersion:       types.StringValue(api.SSLConfig.MinVersion),
			Fingerprint:      types.StringValue(api.SSLConfig.Fingerprint),
			SelfSigned:       types.BoolValue(api.SSLConfig.SelfSigned),
			URL:              types.StringValue(api.SSLConfig.URL),
		}),
	}
	return &model, nil
}

type CheckSSLCertResourceAPI struct {
	provider *providerImpl
}

func (c CheckSSLCertResourceAPI) Create(ctx context.Context, arg upapi.CheckSSLCert) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateSSLCert(ctx, arg)
}

func (c CheckSSLCertResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckSSLCertResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckSSLCert) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateSSLCert(ctx, pk, arg)
}

func (c CheckSSLCertResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
