package provider

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/terraform-provider-uptime/internal/uptimeapi"
)

type tagResourceData struct {
	ID       types.Int64  `tfsdk:"id"            map:"Pk"`
	Tag      types.String `tfsdk:"tag"           map:"Tag"`
	ColorHex types.String `tfsdk:"color_hex"     map:"ColorHex"`
	URL      types.String `tfsdk:"url"           map:"Url"`
}

func (t *tagResourceData) ToAPI(ctx context.Context) (uptimeapi.CheckTag, diag.Diagnostics) {
	obj := uptimeapi.CheckTag{}
	diags := fromTerraform(ctx, &obj, t)
	return obj, diags
}

func (t *tagResourceData) FromAPI(obj uptimeapi.CheckTag) diag.Diagnostics {
	return toTerraform(t, obj)
}

var _ resource.Resource = &tagResourceImpl{}
var _ resource.ResourceWithImportState = &tagResourceImpl{}

type tagResourceImpl struct {
	api uptimeapi.ClientInterface
}

func (r *tagResourceImpl) Metadata(_ context.Context, rq resource.MetadataRequest, rs *resource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_tag"
}

func (r *tagResourceImpl) Schema(_ context.Context, _ resource.SchemaRequest, rs *resource.SchemaResponse) {
	rs.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"tag": schema.StringAttribute{
				Required: true,
			},
			"color_hex": schema.StringAttribute{
				Required: true,
			},
			"url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *tagResourceImpl) Create(ctx context.Context, rq resource.CreateRequest, rs *resource.CreateResponse) {
	data := new(tagResourceData)
	if diag := rq.Config.Get(ctx, &data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	api, diags := data.ToAPI(ctx)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	res, err := r.api.PostServicetaglist(ctx, api)
	if err != nil {
		rs.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", res.Status)
		return
	}

	// Uptime.com API doesn't really follow own OpenAPI spec. Generated client code is not usable. Have to manually
	// redefine partial response JSON structure here.
	obj := struct {
		Result uptimeapi.CheckTag `json:"results"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&obj)
	if err != nil {
		rs.Diagnostics.AddError("Failed to decode response", err.Error())
		return
	}
	diags = data.FromAPI(obj.Result)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = rs.State.Set(ctx, data)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	return
}

func (r *tagResourceImpl) Read(ctx context.Context, rq resource.ReadRequest, rs *resource.ReadResponse) {
	var data tagResourceData
	if diag := rq.State.Get(ctx, &data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	api := uptimeapi.ClientWithResponses{ClientInterface: r.api}
	obj, err := api.GetServiceTagDetailWithResponse(ctx, data.ID.String())
	if err != nil {
		rs.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	if obj.StatusCode() != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", obj.Status())
		return
	}
	data.FromAPI(*obj.JSON200)
	if diag := rs.State.Set(ctx, data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	return
}

func (r *tagResourceImpl) Update(ctx context.Context, rq resource.UpdateRequest, rs *resource.UpdateResponse) {
	prev := new(tagResourceData)
	if diag := rq.State.Get(ctx, prev); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	next := new(tagResourceData)
	if diag := rq.Config.Get(ctx, next); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	api, diags := next.ToAPI(ctx)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	res, err := r.api.PutServiceTagDetail(ctx, prev.ID.String(), api)
	if err != nil {
		rs.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", res.Status)
		return
	}

	// Uptime.com API doesn't really follow own OpenAPI spec. Generated client code is not usable. Have to manually
	// redefine partial response JSON structure here.
	obj := new(struct {
		Result uptimeapi.CheckTag `json:"results"`
	})
	err = json.NewDecoder(res.Body).Decode(obj)
	if err != nil {
		rs.Diagnostics.AddError("Failed to decode response", err.Error())
		return
	}
	prev.FromAPI(obj.Result)

	if res.StatusCode != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", res.Status)
		return
	}
	if diag := rs.State.Set(ctx, prev); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	return
}

func (r *tagResourceImpl) Delete(ctx context.Context, rq resource.DeleteRequest, rs *resource.DeleteResponse) {
	var data tagResourceData
	if diag := rq.State.Get(ctx, &data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	_, err := r.api.DeleteServiceTagDetail(ctx, data.ID.String())
	if err != nil {
		rs.Diagnostics.AddError("Delete failed", err.Error())
		return
	}
}

func (r *tagResourceImpl) ImportState(ctx context.Context, rq resource.ImportStateRequest, rs *resource.ImportStateResponse) {
	api := uptimeapi.ClientWithResponses{ClientInterface: r.api}
	obj, err := api.GetServicetaglistWithResponse(ctx, &uptimeapi.GetServicetaglistParams{
		Search: &rq.ID,
	})
	if err != nil {
		rs.Diagnostics.AddError("Import failed", err.Error())
		return
	}
	if obj.StatusCode() != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", obj.Status())
		return
	}
	for _, tag := range *obj.JSON200.Results {
		if tag.Tag != rq.ID {
			continue
		}
		data := new(tagResourceData)
		data.FromAPI(tag)
		if diag := rs.State.Set(ctx, data); diag.HasError() {
			rs.Diagnostics.Append(diag...)
			return
		}
		return
	}
	rs.Diagnostics.AddError("Import failed", "tag not found")
	return
}
