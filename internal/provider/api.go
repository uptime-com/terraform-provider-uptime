package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

// errResourceGone signals that the API returned a record that still exists at
// the HTTP level but is logically deleted (e.g. carries a deleted_at marker).
// A resource Read may return it, wrapped or not, to have the resource dropped
// from state exactly like a 404 so out-of-band deletions surface as drift.
var errResourceGone = errors.New("resource marked deleted by API")

// isNotFoundError reports whether err represents a resource that no longer
// exists: either an HTTP 404 response from the Uptime.com API or an
// errResourceGone marker. It is used during refresh to detect resources that
// were deleted out-of-band so they can be dropped from state instead of failing.
func isNotFoundError(err error) bool {
	if errors.Is(err, errResourceGone) {
		return true
	}
	var apiErr *upapi.Error
	if errors.As(err, &apiErr) {
		return apiErr.Response != nil && apiErr.Response.StatusCode == http.StatusNotFound
	}
	return false
}

// notFoundWarning surfaces a state removal in the plan output. It fires when the
// API reports a resource as gone (a 404, or a deleted_at marker). A 404 can also
// mean a wrong subaccount/endpoint configuration, in which case silently dropping
// resources would cascade into recreating everything.
func notFoundWarning(typeNameSuffix string, pk upapi.PrimaryKeyable) diag.Diagnostic {
	return diag.NewWarningDiagnostic(
		"Resource Not Found",
		fmt.Sprintf(
			"uptime_%s with ID %d no longer exists on the server and was removed "+
				"from state. If the resource was not deleted out-of-band, check the "+
				"provider's subaccount and endpoint configuration before applying.",
			typeNameSuffix, pk.PrimaryKey(),
		),
	)
}

type API[A, R any] interface {
	Create(context.Context, A) (*R, error)
	Read(context.Context, upapi.PrimaryKeyable) (*R, error)
	Update(context.Context, upapi.PrimaryKeyable, A) (*R, error)
	Delete(context.Context, upapi.PrimaryKeyable) error
}

type APIModel interface {
	upapi.PrimaryKeyable
}

type StateGetter interface {
	Get(context.Context, interface{}) diag.Diagnostics
}

type APIModeler[M APIModel, A, R any] interface {
	Get(context.Context, StateGetter) (*M, diag.Diagnostics)
	ToAPIArgument(M) (*A, error)
	FromAPIResult(R) (*M, error)
}

// PlanValuePreserver is an optional interface that APIModelers can implement
// to preserve plan values for fields that the API does not return. It is applied
// on Create and Update (apply time), and on Read unless the modeler also
// implements ReadValuePreserver.
type PlanValuePreserver[M APIModel] interface {
	PreservePlanValues(result *M, plan *M) *M
}

// ReadValuePreserver is an optional interface that APIModelers can implement to
// override how prior state is reconciled with the API response on Read only.
// When implemented, Read calls PreserveReadValues instead of PreservePlanValues,
// letting a resource trust the server response on refresh so out-of-band changes
// surface as drift.
type ReadValuePreserver[M APIModel] interface {
	PreserveReadValues(result *M, state *M) *M
}

type APIResourceMetadata struct {
	schema.Schema
	TypeNameSuffix   string
	ConfigValidators func(context.Context) []resource.ConfigValidator
}

type APIResource[M APIModel, A, R any] struct {
	api  API[A, R]
	mod  APIModeler[M, A, R]
	meta APIResourceMetadata
}

func (r APIResource[M, A, R]) Metadata(_ context.Context, rq resource.MetadataRequest, rs *resource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_" + r.meta.TypeNameSuffix
}

func (r APIResource[M, A, R]) Schema(_ context.Context, _ resource.SchemaRequest, rs *resource.SchemaResponse) {
	rs.Schema = r.meta.Schema
}

func (r APIResource[M, A, R]) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if r.meta.ConfigValidators == nil {
		return []resource.ConfigValidator{}
	}
	return r.meta.ConfigValidators(ctx)
}

const (
	toAPIArgumentError = "To API Argument Conversion Failed"
	fromAPIResultError = "From API Result Conversion Failed"

	apiOperationCreate = "API Create Operation Failed"
	apiOperationRead   = "API Read Operation Failed"
	apiOperationUpdate = "API Update Operation Failed"
	apiOperationDelete = "API Delete Operation Failed"
)

func (r APIResource[M, A, R]) apiOperationError(op string, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(op, err.Error())
}

func (r APIResource[M, A, R]) apiConversionError(op string, src, dst any, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(op,
		fmt.Sprintf("%T -> %T: %s", src, dst, err.Error()),
	)
}

func (r APIResource[M, A, R]) Create(ctx context.Context, rq resource.CreateRequest, rs *resource.CreateResponse) {
	planModel, diags := r.mod.Get(ctx, rq.Plan)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	arg, err := r.mod.ToAPIArgument(*planModel)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(toAPIArgumentError, planModel, arg, err))
		return
	}

	res, err := r.api.Create(ctx, *arg)
	if err != nil {
		rs.Diagnostics.Append(r.apiOperationError(apiOperationCreate, err))
		return
	}

	resultModel, err := r.mod.FromAPIResult(*res)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(fromAPIResultError, res, resultModel, err))
		return
	}

	// If the modeler implements PlanValuePreserver, use it to preserve plan values
	// for fields that the API doesn't return (like sensitive fields)
	if preserver, ok := any(r.mod).(PlanValuePreserver[M]); ok {
		resultModel = preserver.PreservePlanValues(resultModel, planModel)
	}

	diags = rs.State.Set(ctx, resultModel)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}
	return
}

func (r APIResource[M, A, R]) Read(ctx context.Context, rq resource.ReadRequest, rs *resource.ReadResponse) {
	stateModel, diags := r.mod.Get(ctx, rq.State)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	res, err := r.api.Read(ctx, *stateModel)
	if err != nil {
		if isNotFoundError(err) {
			rs.Diagnostics.Append(notFoundWarning(r.meta.TypeNameSuffix, (*stateModel).PrimaryKey()))
			rs.State.RemoveResource(ctx)
			return
		}
		rs.Diagnostics.Append(r.apiOperationError(apiOperationRead, err))
		return
	}

	resultModel, err := r.mod.FromAPIResult(*res)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(fromAPIResultError, res, resultModel, err))
		return
	}

	// On refresh, prefer a Read-specific reconciler when the modeler provides one
	// (it may trust the API to surface out-of-band drift); otherwise fall back to
	// the apply-time plan-value preservation for fields the API doesn't return.
	if preserver, ok := any(r.mod).(ReadValuePreserver[M]); ok {
		resultModel = preserver.PreserveReadValues(resultModel, stateModel)
	} else if preserver, ok := any(r.mod).(PlanValuePreserver[M]); ok {
		resultModel = preserver.PreservePlanValues(resultModel, stateModel)
	}

	diags = rs.State.Set(ctx, resultModel)
	if rs.Diagnostics.HasError() {
		return
	}
	return
}

func (r APIResource[M, A, R]) Update(ctx context.Context, rq resource.UpdateRequest, rs *resource.UpdateResponse) {
	state, diags := r.mod.Get(ctx, rq.State)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	planModel, diags := r.mod.Get(ctx, rq.Plan)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	arg, err := r.mod.ToAPIArgument(*planModel)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(toAPIArgumentError, planModel, arg, err))
		return
	}

	res, err := r.api.Update(ctx, *state, *arg)
	if err != nil {
		rs.Diagnostics.Append(r.apiOperationError(apiOperationUpdate, err))
		return
	}

	resultModel, err := r.mod.FromAPIResult(*res)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(fromAPIResultError, res, resultModel, err))
		return
	}

	// Preserve plan values for fields that the API doesn't return
	if preserver, ok := any(r.mod).(PlanValuePreserver[M]); ok {
		resultModel = preserver.PreservePlanValues(resultModel, planModel)
	}

	diags = rs.State.Set(ctx, resultModel)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}
	return
}

func (r APIResource[M, A, R]) Delete(ctx context.Context, rq resource.DeleteRequest, rs *resource.DeleteResponse) {
	state, diags := r.mod.Get(ctx, rq.State)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	if err := r.api.Delete(ctx, *state); err != nil {
		rs.Diagnostics.Append(r.apiOperationError(apiOperationDelete, err))
		return
	}

	return
}

// ImportableAPIResource wraps APIResource and adds import support.
// Use this for resources that need import functionality.
type ImportableAPIResource[M APIModel, A, R any] struct {
	APIResource[M, A, R]
	importHandler func(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse)
}

// NewImportableAPIResource creates a new ImportableAPIResource with the given import handler.
func NewImportableAPIResource[M APIModel, A, R any](
	api API[A, R],
	mod APIModeler[M, A, R],
	meta APIResourceMetadata,
	importHandler func(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse),
) ImportableAPIResource[M, A, R] {
	return ImportableAPIResource[M, A, R]{
		APIResource: APIResource[M, A, R]{
			api:  api,
			mod:  mod,
			meta: meta,
		},
		importHandler: importHandler,
	}
}

// ImportState implements resource.ResourceWithImportState
func (r ImportableAPIResource[M, A, R]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.importHandler(ctx, req, resp)
}
