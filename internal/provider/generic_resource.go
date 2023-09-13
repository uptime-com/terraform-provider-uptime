package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

type genericResourceMetadata struct {
	TypeNameSuffix string
	Schema         schema.Schema
}

type genericResource[Model, Arg, Res any] struct {
	api      genericResourceAPI[Arg, Res]
	metadata genericResourceMetadata
}

type genericResourceAPI[Arg, Res any] interface {
	Create(context.Context, Arg) (*Res, error)
	Read(context.Context, upapi.PrimaryKeyable) (*Res, error)
	Update(context.Context, upapi.PrimaryKeyable, Arg) (*Res, error)
	Delete(context.Context, upapi.PrimaryKeyable) error
}

func (r *genericResource[Model, Arg, Res]) Metadata(_ context.Context, rq resource.MetadataRequest, rs *resource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_" + r.metadata.TypeNameSuffix
}

func (r *genericResource[Model, Arg, Res]) Schema(_ context.Context, rq resource.SchemaRequest, rs *resource.SchemaResponse) {
	rs.Schema = r.metadata.Schema
}

func (r *genericResource[Model, Arg, Res]) Create(ctx context.Context, rq resource.CreateRequest, rs *resource.CreateResponse) {
	var (
		obj    Model
		apiarg Arg
		apires *Res
		diags  diag.Diagnostics
		err    error
	)
	diags = rq.Plan.Get(ctx, &obj)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = valueToAPI(&apiarg, obj)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	apires, err = r.api.Create(ctx, apiarg)
	if err != nil {
		rs.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	diags = valueFromAPI(&obj, *apires)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = rs.State.Set(ctx, &obj)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	return
}

func (r *genericResource[Model, Arg, Res]) Read(ctx context.Context, rq resource.ReadRequest, rs *resource.ReadResponse) {
	var (
		pk     upapi.PrimaryKeyable
		obj    Model
		apires *Res
		diags  diag.Diagnostics
		err    error
	)
	pk, diags = primaryKeyFromState(ctx, rq.State)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	apires, err = r.api.Read(ctx, pk)
	if err != nil {
		rs.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	diags = valueFromAPI(&obj, *apires)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = rs.State.Set(ctx, &obj)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	return
}

func (r *genericResource[Model, Arg, Res]) Update(ctx context.Context, rq resource.UpdateRequest, rs *resource.UpdateResponse) {
	var (
		pk     upapi.PrimaryKeyable
		obj    Model
		apiarg Arg
		apires *Res
	)
	pk, diags := primaryKeyFromState(ctx, rq.State)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = rq.Plan.Get(ctx, &obj)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = valueToAPI(&apiarg, obj)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	apires, err := r.api.Update(ctx, pk, apiarg)
	if err != nil {
		rs.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	diags = valueFromAPI(&obj, *apires)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = rs.State.Set(ctx, &obj)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	return
}

func (r *genericResource[Model, Arg, Res]) Delete(ctx context.Context, rq resource.DeleteRequest, rs *resource.DeleteResponse) {
	var (
		pk    upapi.PrimaryKeyable
		diags diag.Diagnostics
		err   error
	)
	pk, diags = primaryKeyFromState(ctx, rq.State)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	err = r.api.Delete(ctx, pk)
	if err != nil {
		rs.Diagnostics.AddError("Delete failed", err.Error())
		return
	}
	return
}

func primaryKeyFromState(ctx context.Context, state tfsdk.State) (upapi.PrimaryKeyable, diag.Diagnostics) {
	var pk types.Int64
	diags := state.GetAttribute(ctx, path.Root("id"), &pk)
	if diags.HasError() {
		return nil, diags
	}
	return upapi.PrimaryKey(pk.ValueInt64()), nil
}
