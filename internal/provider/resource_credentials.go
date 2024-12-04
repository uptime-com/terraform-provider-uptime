package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCredentialResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CredentialResourceModel, upapi.Credential, upapi.Credential]{
		api: CredentialResourceAPI{provider: p},
		mod: CredentialResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "credential",
			Schema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": IDSchemaAttribute(),
					"display_name": schema.StringAttribute{
						Required: true,
					},
					"description": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"credential_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							OneOfStringValidator([]string{"BASIC", "CERTIFICATE", "TOKEN"}),
						},
					},
					"username": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"secret": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"certificate": schema.StringAttribute{
								Computed:  true,
								Optional:  true,
								Sensitive: true,
							},
							"key": schema.StringAttribute{
								Computed:  true,
								Optional:  true,
								Sensitive: true,
							},
							"password": schema.StringAttribute{
								Computed:  true,
								Optional:  true,
								Sensitive: true,
							},
							"passphrase": schema.StringAttribute{
								Computed:  true,
								Optional:  true,
								Sensitive: true,
							},
							"secret": schema.StringAttribute{
								Computed:  true,
								Optional:  true,
								Sensitive: true,
							},
						},
					},
				},
			},
			ConfigValidators: func(context.Context) []resource.ConfigValidator {
				return []resource.ConfigValidator{NewCredentialTypeValidator()}
			},
		},
	}
}

type CredentialResourceModel struct {
	ID             types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	CredentialType types.String `tfsdk:"credential_type"`
	Username       types.String `tfsdk:"username"`
	Secret         types.Object `tfsdk:"secret"`

	secret *CredentialSecretAttribute
}

func (m CredentialResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CredentialSecretAttribute struct {
	Certificate types.String `tfsdk:"certificate"`
	Key         types.String `tfsdk:"key"`
	Password    types.String `tfsdk:"password"`
	Passphrase  types.String `tfsdk:"passphrase"`
	Secret      types.String `tfsdk:"secret"`
}

type CredentialResourceModelAdapter struct {
	SetAttributeAdapter[string]
}

func (a CredentialResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CredentialResourceModel, diag.Diagnostics) {
	model := *new(CredentialResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.secret, diags = a.SecretAttributeContext(ctx, model.Secret)
	if diags.HasError() {
		return nil, diags
	}

	return &model, nil
}

func (a CredentialResourceModelAdapter) ToAPIArgument(model CredentialResourceModel) (*upapi.Credential, error) {
	api := upapi.Credential{
		PK:             model.ID.ValueInt64(),
		DisplayName:    model.DisplayName.ValueString(),
		Description:    model.Description.ValueString(),
		CredentialType: model.CredentialType.ValueString(),
		Username:       model.Username.ValueString(),
		Secret: upapi.CredentialSecret{
			Certificate: model.secret.Certificate.ValueString(),
			Key:         model.secret.Key.ValueString(),
			Password:    model.secret.Password.ValueString(),
			Passphrase:  model.secret.Passphrase.ValueString(),
			Secret:      model.secret.Secret.ValueString(),
		},
	}
	return &api, nil
}

func (a CredentialResourceModelAdapter) FromAPIResult(api upapi.Credential) (*CredentialResourceModel, error) {
	model := CredentialResourceModel{
		ID:             types.Int64Value(api.PK),
		DisplayName:    types.StringValue(api.DisplayName),
		Description:    types.StringValue(api.Description),
		CredentialType: types.StringValue(api.CredentialType),
		Username:       types.StringValue(api.Username),
		Secret: a.SecretAttributeValue(CredentialSecretAttribute{
			Certificate: types.StringValue(api.Secret.Certificate),
			Key:         types.StringValue(api.Secret.Key),
			Passphrase:  types.StringValue(api.Secret.Passphrase),
			Password:    types.StringValue(api.Secret.Password),
			Secret:      types.StringValue(api.Secret.Secret),
		}),
	}
	return &model, nil
}

func (a CredentialResourceModelAdapter) secretAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"certificate": types.StringType,
		"key":         types.StringType,
		"password":    types.StringType,
		"passphrase":  types.StringType,
		"secret":      types.StringType,
	}
}

func (a CredentialResourceModelAdapter) secretAttributeValues(m CredentialSecretAttribute) map[string]attr.Value {
	return map[string]attr.Value{
		"certificate": m.Certificate,
		"key":         m.Key,
		"password":    m.Password,
		"passphrase":  m.Passphrase,
		"secret":      m.Secret,
	}
}

func (a CredentialResourceModelAdapter) SecretAttributeContext(ctx context.Context, v types.Object) (*CredentialSecretAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := new(CredentialSecretAttribute)
	diags := v.As(ctx, m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return m, nil
}

func (a CredentialResourceModelAdapter) SecretAttributeValue(m CredentialSecretAttribute) types.Object {
	return types.ObjectValueMust(a.secretAttributeTypes(), a.secretAttributeValues(m))
}

type CredentialResourceAPI struct {
	provider *providerImpl
}

func (c CredentialResourceAPI) Create(ctx context.Context, arg upapi.Credential) (*upapi.Credential, error) {
	obj, err := c.provider.api.Credentials().Create(ctx, arg)
	obj.Secret = arg.Secret
	return obj, err
}

func (c CredentialResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Credential, error) {
	obj, err := c.provider.api.Credentials().Get(ctx, pk)
	secret := pk.(CredentialResourceModel).secret
	obj.Secret = upapi.CredentialSecret{
		Certificate: secret.Certificate.ValueString(),
		Key:         secret.Key.ValueString(),
		Passphrase:  secret.Passphrase.ValueString(),
		Password:    secret.Password.ValueString(),
		Secret:      secret.Secret.ValueString(),
	}
	return obj, err
}

func (c CredentialResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.Credential) (*upapi.Credential, error) {
	if err := c.Delete(ctx, pk); err != nil {
		return nil, err
	}
	return c.Create(ctx, arg)
}

func (c CredentialResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Credentials().Delete(ctx, pk)
}

type credentialTypeValidator struct{}

func NewCredentialTypeValidator() resource.ConfigValidator {
	return &credentialTypeValidator{}
}

func (v *credentialTypeValidator) Description(ctx context.Context) string {
	return "Validates that the credential_type field has valid values and corresponding secret fields are set correctly."
}

func (v *credentialTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *credentialTypeValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model CredentialResourceModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if model.CredentialType.IsNull() || model.CredentialType.IsUnknown() {
		return
	}

	// var a attr.Value
	var secretAttr CredentialSecretAttribute
	p := path.Root("secret")
	diags = req.Config.GetAttribute(ctx, p, &secretAttr)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	switch model.CredentialType.ValueString() {
	case "BASIC":
		if secretAttr.Password.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"When credential_type is BASIC, the password field must be set.",
			)
		}
		if !secretAttr.Certificate.IsNull() || !secretAttr.Key.IsNull() || !secretAttr.Passphrase.IsNull() || !secretAttr.Secret.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"When credential_type is BASIC, only the password field should be set.",
			)
		}
	case "CERTIFICATE":
		if secretAttr.Certificate.IsNull() || secretAttr.Key.IsNull() || secretAttr.Passphrase.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"When credential_type is CERTIFICATE, the certificate, key, and passphrase fields must be set.",
			)
		}
		if !secretAttr.Password.IsNull() || !secretAttr.Secret.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"When credential_type is CERTIFICATE, only the certificate, key, and passphrase fields should be set.",
			)
		}
	case "TOKEN":
		if model.Secret.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"When credential_type is TOKEN, the secret field must be set.",
			)
		}
		if !secretAttr.Password.IsNull() || !secretAttr.Certificate.IsNull() || !secretAttr.Key.IsNull() || !secretAttr.Passphrase.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				"When credential_type is TOKEN, only the secret field should be set.",
			)
		}
	default:
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			fmt.Sprintf("Invalid credential_type: %s", model.CredentialType.ValueString()),
		)
	}
}
