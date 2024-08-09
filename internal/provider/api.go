package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

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
	return diag.NewErrorDiagnostic(op, fmt.Sprintf(err.Error()))
}

func (r APIResource[M, A, R]) apiConversionError(op string, src, dst any, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(op,
		fmt.Sprintf("%T -> %T: %s", src, dst, err.Error()),
	)
}

func (r APIResource[M, A, R]) Create(ctx context.Context, rq resource.CreateRequest, rs *resource.CreateResponse) {
	model, diags := r.mod.Get(ctx, rq.Plan)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	arg, err := r.mod.ToAPIArgument(*model)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(toAPIArgumentError, model, arg, err))
		return
	}

	res, err := r.api.Create(ctx, *arg)
	if err != nil {
		rs.Diagnostics.Append(r.apiOperationError(apiOperationCreate, err))
		return
	}

	model, err = r.mod.FromAPIResult(*res)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(fromAPIResultError, res, model, err))
		return
	}

	diags = rs.State.Set(ctx, model)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}
	return
}

func (r APIResource[M, A, R]) Read(ctx context.Context, rq resource.ReadRequest, rs *resource.ReadResponse) {
	model, diags := r.mod.Get(ctx, rq.State)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	res, err := r.api.Read(ctx, *model)
	if err != nil {
		rs.Diagnostics.Append(r.apiOperationError(apiOperationRead, err))
		return
	}

	model, err = r.mod.FromAPIResult(*res)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(fromAPIResultError, res, model, err))
		return
	}

	diags = rs.State.Set(ctx, model)
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

	plan, diags := r.mod.Get(ctx, rq.Plan)
	rs.Diagnostics.Append(diags...)
	if rs.Diagnostics.HasError() {
		return
	}

	arg, err := r.mod.ToAPIArgument(*plan)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(toAPIArgumentError, plan, arg, err))
		return
	}

	res, err := r.api.Update(ctx, *state, *arg)
	if err != nil {
		rs.Diagnostics.Append(r.apiOperationError(apiOperationUpdate, err))
		return
	}

	state, err = r.mod.FromAPIResult(*res)
	if err != nil {
		rs.Diagnostics.Append(r.apiConversionError(fromAPIResultError, res, state, err))
		return
	}

	diags = rs.State.Set(ctx, state)
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
