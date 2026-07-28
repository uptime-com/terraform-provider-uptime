package provider

import (
	"context"
	"errors"
	"fmt"

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
//
// Import uses the composite ID "service_id:variable_id"; see importServiceVariableState.
func NewServiceVariableResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[ServiceVariableResourceModel, ServiceVariableWrapper, ServiceVariableWrapper](
		ServiceVariableResourceAPI{provider: p},
		ServiceVariableResourceModelAdapter{},
		APIResourceMetadata{
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
		importServiceVariableState(p),
	)
}

// importServiceVariableState imports a service variable by "service_id:variable_id".
//
// The composite form is required because the API never returns service_id: Read sources
// it from state, and it is Required + RequiresReplace, so importing by variable ID alone
// would leave it at zero and make the very next plan propose a replacement.
//
// Nothing server-side ties the two halves together - the GET is keyed by variable ID
// alone - so a mistyped service_id would otherwise be accepted and only surface later as
// that same silent replacement, destroying and recreating the link. Two extra reads at
// import time buy that check: the variable names its owning check in `service`, and the
// check named by service_id reports its own name, so a mismatch is caught up front.
// Import is a one-off operation, so the extra requests do not affect plan or apply.
func importServiceVariableState(
	p *providerImpl,
) func(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		serviceID, variableID, err := ParseCompositeID(req.ID, "service_id")
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", err.Error())
			return
		}

		variable, err := p.api.ServiceVariables().Get(ctx, upapi.PrimaryKey(variableID))
		if err != nil {
			resp.Diagnostics.AddError("Service Variable Not Found",
				fmt.Sprintf("could not read service variable %d: %s", variableID, err.Error()))
			return
		}
		check, err := p.api.Checks().Get(ctx, upapi.PrimaryKey(serviceID))
		if err != nil {
			resp.Diagnostics.AddError("Check Not Found",
				fmt.Sprintf("could not read check %d named by service_id: %s", serviceID, err.Error()))
			return
		}

		// Only a positive mismatch is rejected: `service` is omitempty on both sides, and
		// refusing an import because a name came back blank would be worse than importing.
		if variable.Service != "" && check.Name != "" && variable.Service != check.Name {
			resp.Diagnostics.AddError("Import ID Mismatch", fmt.Sprintf(
				"service variable %d belongs to check %q, but service_id %d is check %q. "+
					"Use the ID of the check that owns the variable.",
				variableID, variable.Service, serviceID, check.Name))
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), serviceID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), variableID)...)
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

// errServiceVariableNoID reports a write the endpoint rejected with HTTP 200 and an
// empty results object, which decodes to a zero-valued record. Because id is Computed,
// persisting it would make Terraform accept id 0 as a valid apply result.
var errServiceVariableNoID = errors.New(
	"service variable write returned no ID: the API rejected the request without reporting an error status. " +
		"A variable with this name may already exist on the check; verify variable_name, credential_id, " +
		"service_id and the configured subaccount",
)

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
	if result == nil || result.ID == 0 {
		return nil, errServiceVariableNoID
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
	// service_id is not part of the API response, so it can only come from state - which
	// is why import must carry it. ServiceVariableWrapper also implements PrimaryKeyable,
	// so assert rather than panic if a caller ever passes the wrong one.
	model, ok := pk.(ServiceVariableResourceModel)
	if !ok {
		return nil, fmt.Errorf("service variable read: expected %T, got %T", model, pk)
	}
	result, err := c.provider.api.ServiceVariables().Get(ctx, pk)
	if err != nil {
		return nil, err
	}
	// A link removed from the check in the UI is hard-deleted, so a later Get
	// answers HTTP 200 with an empty body that decodes to a zero-valued record
	// (ID 0, no deleted_at). A soft-deleted link instead comes back flagged with
	// deleted_at. Treat either as gone so the removal surfaces as drift and the
	// next apply recreates the link, rather than being masked by the apply-time
	// backfill or issuing a no-op update against ID 0 that never restores it.
	if result.ID == 0 || result.DeletedAt != nil {
		return nil, errResourceGone
	}
	// Extract credential_id from nested credential object if not at top level
	if result.CredentialID == 0 && result.Credential != nil {
		result.CredentialID = result.Credential.ID
	}
	return &ServiceVariableWrapper{
		ServiceVariable: *result,
		ServiceID:       model.ServiceID.ValueInt64(),
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
	if result == nil || result.ID == 0 {
		return nil, errServiceVariableNoID
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
