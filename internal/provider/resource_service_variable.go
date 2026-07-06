package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

// NewServiceVariableResource creates a new service variable resource.
//
// Service variables allow you to securely inject credential properties into check configurations
// without exposing sensitive values. This is useful for authentication, API tokens, certificates,
// and other sensitive data that checks need to access.
func NewServiceVariableResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[ServiceVariableResourceModel, ServiceVariableWrapper, ServiceVariableWrapper]{
		api: ServiceVariableResourceAPI{provider: p},
		mod: ServiceVariableResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "service_variable",
			Schema: schema.Schema{
				Description: "Links a credential property to a check/service, allowing secure injection of sensitive values into check configurations.",
				Attributes: map[string]schema.Attribute{
					"id": IDSchemaAttribute(),
					"service_id": schema.Int64Attribute{
						Required: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
						Description: "The ID of the check/service this variable belongs to. Changing this forces recreation of the resource.",
					},
					"credential_id": schema.Int64Attribute{
						Required:    true,
						Description: "The ID of the credential containing the sensitive value to inject",
					},
					"variable_name": schema.StringAttribute{
						Required:    true,
						Description: "The name of the variable as referenced in the check configuration",
					},
					"property_name": schema.StringAttribute{
						Required:    true,
						Description: "The property name from the credential to use. Valid values depend on credential_type: 'password' for BASIC, 'secret' for TOKEN, 'certificate'/'key'/'passphrase' for CERTIFICATE",
					},
					"service": schema.StringAttribute{
						Computed:    true,
						Description: "Service identifier (computed)",
					},
					"account": schema.StringAttribute{
						Computed:    true,
						Description: "Account identifier (computed)",
					},
				},
			},
		},
	}
}

type ServiceVariableWrapper struct {
	upapi.ServiceVariable
	ServiceID int64
}

func (w ServiceVariableWrapper) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(w.ID)
}

type ServiceVariableResourceModel struct {
	ID           types.Int64  `tfsdk:"id" ref:"PK,opt"`
	ServiceID    types.Int64  `tfsdk:"service_id"`
	CredentialID types.Int64  `tfsdk:"credential_id"`
	VariableName types.String `tfsdk:"variable_name"`
	PropertyName types.String `tfsdk:"property_name"`
	Service      types.String `tfsdk:"service"`
	Account      types.String `tfsdk:"account"`
}

func (m ServiceVariableResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type ServiceVariableResourceModelAdapter struct {
	SetAttributeAdapter[int64]
}

func (a ServiceVariableResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*ServiceVariableResourceModel, diag.Diagnostics) {
	model := *new(ServiceVariableResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a ServiceVariableResourceModelAdapter) ToAPIArgument(model ServiceVariableResourceModel) (*ServiceVariableWrapper, error) {
	return &ServiceVariableWrapper{
		ServiceID: model.ServiceID.ValueInt64(),
		ServiceVariable: upapi.ServiceVariable{
			// ID should not be set for create/update - it comes from the API response
			CredentialID: model.CredentialID.ValueInt64(),
			VariableName: model.VariableName.ValueString(),
			PropertyName: model.PropertyName.ValueString(),
		},
	}, nil
}

func (a ServiceVariableResourceModelAdapter) FromAPIResult(api ServiceVariableWrapper) (*ServiceVariableResourceModel, error) {
	return &ServiceVariableResourceModel{
		ID:           types.Int64Value(api.ID),
		ServiceID:    types.Int64Value(api.ServiceID), // Preserved from wrapper
		CredentialID: types.Int64Value(api.CredentialID),
		VariableName: types.StringValue(api.VariableName),
		PropertyName: types.StringValue(api.PropertyName),
		Service:      types.StringValue(api.Service),
		Account:      types.StringValue(api.Account),
	}, nil
}

// PreservePlanValues restores credential_id, property_name and variable_name from
// the plan when the API response omits them. These are Required attributes, so an
// empty 0/"" result mismatches the planned value and trips "Provider produced
// inconsistent result after apply". This runs at apply time only (Create/Update);
// PreserveReadValues deliberately opts refresh out of this backfill so that
// out-of-band deletions are not masked.
func (a ServiceVariableResourceModelAdapter) PreservePlanValues(result, plan *ServiceVariableResourceModel) *ServiceVariableResourceModel {
	if result.CredentialID.ValueInt64() == 0 {
		result.CredentialID = plan.CredentialID
	}
	if result.PropertyName.ValueString() == "" {
		result.PropertyName = plan.PropertyName
	}
	if result.VariableName.ValueString() == "" {
		result.VariableName = plan.VariableName
	}
	return result
}

// PreserveReadValues trusts the API response on refresh: it does not backfill from
// prior state. Combined with the deleted_at handling in Read, this lets a UI-side
// removal of the credential link show up as drift instead of being silently masked.
func (a ServiceVariableResourceModelAdapter) PreserveReadValues(result, _ *ServiceVariableResourceModel) *ServiceVariableResourceModel {
	return result
}

type ServiceVariableResourceAPI struct {
	provider *providerImpl
}

func (c ServiceVariableResourceAPI) Create(ctx context.Context, arg ServiceVariableWrapper) (*ServiceVariableWrapper, error) {
	createReq := upapi.ServiceVariableCreateRequest{
		ServiceID:    arg.ServiceID,
		CredentialID: arg.CredentialID,
		VariableName: arg.VariableName,
		PropertyName: arg.PropertyName,
	}
	result, err := c.provider.api.ServiceVariables().Create(ctx, createReq)
	if err != nil {
		return nil, err
	}
	// Extract credential_id from nested credential object if not at top level
	if result.CredentialID == 0 && result.Credential != nil {
		result.CredentialID = result.Credential.ID
	}
	return &ServiceVariableWrapper{
		ServiceVariable: *result,
		ServiceID:       arg.ServiceID,
	}, nil
}

func (c ServiceVariableResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*ServiceVariableWrapper, error) {
	// Extract ServiceID from the wrapper
	wrapper := pk.(ServiceVariableResourceModel)
	result, err := c.provider.api.ServiceVariables().Get(ctx, pk)
	if err != nil {
		return nil, err
	}
	// The API keeps returning the link after it is unbound in the UI, only
	// flagging it with deleted_at. Treat that as gone so the deletion surfaces
	// as drift instead of being masked by the apply-time backfill.
	if result.DeletedAt != nil {
		return nil, errResourceGone
	}
	// Extract credential_id from nested credential object if not at top level
	if result.CredentialID == 0 && result.Credential != nil {
		result.CredentialID = result.Credential.ID
	}
	return &ServiceVariableWrapper{
		ServiceVariable: *result,
		ServiceID:       wrapper.ServiceID.ValueInt64(),
	}, nil
}

func (c ServiceVariableResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg ServiceVariableWrapper) (*ServiceVariableWrapper, error) {
	updateReq := upapi.ServiceVariableUpdateRequest{
		ServiceID:    arg.ServiceID,
		CredentialID: arg.CredentialID,
		VariableName: arg.VariableName,
		PropertyName: arg.PropertyName,
	}
	result, err := c.provider.api.ServiceVariables().Update(ctx, pk, updateReq)
	if err != nil {
		return nil, err
	}
	// Extract credential_id from nested credential object if not at top level
	if result.CredentialID == 0 && result.Credential != nil {
		result.CredentialID = result.Credential.ID
	}
	return &ServiceVariableWrapper{
		ServiceVariable: *result,
		ServiceID:       arg.ServiceID,
	}, nil
}

func (c ServiceVariableResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.ServiceVariables().Delete(ctx, pk)
}

func (c ServiceVariableResourceAPI) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
