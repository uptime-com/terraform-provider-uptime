package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCredentialsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return CredentialsDataSource{p: p}
}

// CredentialsDataSchema defines the schema for the credentials data source.
var CredentialsDataSchema = schema.Schema{
	Description: "Retrieve a list of all credentials configured in your Uptime.com account. Credentials can be used for authentication in various check types.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"credentials": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all credentials in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the credential",
					},
					"display_name": schema.StringAttribute{
						Computed:    true,
						Description: "Human-readable name for the credential",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Optional description providing additional context about the credential",
					},
					"credential_type": schema.StringAttribute{
						Computed:    true,
						Description: "Type of credential. Valid values: BASIC (username/password), CERTIFICATE (SSL/TLS certificate), TOKEN (API token or secret)",
					},
					"hint": schema.StringAttribute{
						Computed:    true,
						Description: "A hint or reminder about the credential (e.g., partial credential value or usage note)",
					},
					"username": schema.StringAttribute{
						Computed:    true,
						Description: "Username for BASIC authentication credentials. Empty for other credential types",
					},
					"version": schema.StringAttribute{
						Computed:    true,
						Description: "Version information for the credential",
					},
					"used_secret_properties": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of secret property names that are populated for this credential (e.g., ['password'], ['certificate', 'key', 'passphrase'], or ['secret'])",
					},
					"created_by": schema.Int64Attribute{
						Computed:    true,
						Description: "User ID of the account member who created this credential",
					},
				},
			},
		},
	},
}

type CredentialsDataSourceModel struct {
	ID          types.String                     `tfsdk:"id"`
	Credentials []CredentialsDataSourceItemModel `tfsdk:"credentials"`
}

type CredentialsDataSourceItemModel struct {
	ID                   types.Int64  `tfsdk:"id"`
	DisplayName          types.String `tfsdk:"display_name"`
	Description          types.String `tfsdk:"description"`
	CredentialType       types.String `tfsdk:"credential_type"`
	Hint                 types.String `tfsdk:"hint"`
	Username             types.String `tfsdk:"username"`
	Version              types.String `tfsdk:"version"`
	UsedSecretProperties types.List   `tfsdk:"used_secret_properties"`
	CreatedBy            types.Int64  `tfsdk:"created_by"`
}

var _ datasource.DataSource = &CredentialsDataSource{}

type CredentialsDataSource struct {
	p *providerImpl
}

func (d CredentialsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_credentials"
}

func (d CredentialsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = CredentialsDataSchema
}

func (d CredentialsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.Credentials().List(ctx, upapi.CredentialListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	model := CredentialsDataSourceModel{
		ID:          types.StringValue(""),
		Credentials: make([]CredentialsDataSourceItemModel, len(api)),
	}

	for i := range api {
		// Convert UsedSecretProperties to types.List
		var usedSecretPropsList types.List
		if len(api[i].UsedSecretProperties) > 0 {
			elements := make([]attr.Value, len(api[i].UsedSecretProperties))
			for j, prop := range api[i].UsedSecretProperties {
				elements[j] = types.StringValue(prop)
			}
			usedSecretPropsList = types.ListValueMust(types.StringType, elements)
		} else {
			usedSecretPropsList = types.ListNull(types.StringType)
		}

		model.Credentials[i] = CredentialsDataSourceItemModel{
			ID:                   types.Int64Value(api[i].PK),
			DisplayName:          types.StringValue(api[i].DisplayName),
			Description:          types.StringValue(api[i].Description),
			CredentialType:       types.StringValue(api[i].CredentialType),
			Hint:                 types.StringValue(api[i].Hint),
			Username:             types.StringValue(api[i].Username),
			Version:              types.StringValue(api[i].Version),
			UsedSecretProperties: usedSecretPropsList,
			CreatedBy:            types.Int64Value(api[i].CreatedBy),
		}
	}

	rs.Diagnostics = rs.State.Set(ctx, model)
}
