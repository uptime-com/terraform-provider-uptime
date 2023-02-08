package provider

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/uptime-com/terraform-provider-uptime/internal/uptimeapi"
)

type tagResourceData struct {
	ID       *string `tfsdk:"id"            api:"Pk"`
	Tag      string  `tfsdk:"tag"`
	ColorHex string  `tfsdk:"color_hex"`
	URL      *string `tfsdk:"url"           api:"Url"`
}

func (t *tagResourceData) ToAPI() uptimeapi.CheckTag {
	obj := uptimeapi.CheckTag{}
	mirror(&obj, t)
	return obj
}

func (t *tagResourceData) FromAPI(obj uptimeapi.CheckTag) {
	mirror(t, obj)
}

var _ resource.Resource = &tagResourceImpl{}
var _ resource.ResourceWithImportState = &tagResourceImpl{}

type tagResourceImpl struct {
	api uptimeapi.ClientInterface
}

func (t *tagResourceImpl) Metadata(_ context.Context, rq resource.MetadataRequest, rs *resource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_tag"
}

func (t *tagResourceImpl) Schema(_ context.Context, _ resource.SchemaRequest, rs *resource.SchemaResponse) {
	rs.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
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

func (t *tagResourceImpl) Create(ctx context.Context, rq resource.CreateRequest, rs *resource.CreateResponse) {
	data := new(tagResourceData)
	if diag := rq.Config.Get(ctx, &data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	res, err := t.api.PostServicetaglist(ctx, data.ToAPI())
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
	var obj struct {
		Result uptimeapi.CheckTag `json:"results"`
	}
	err = json.NewDecoder(res.Body).Decode(&obj)
	if err != nil {
		rs.Diagnostics.AddError("Failed to decode response", err.Error())
		return
	}
	data.FromAPI(obj.Result)
	if diags := rs.State.Set(ctx, data); diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	return
}

func (t *tagResourceImpl) Read(ctx context.Context, rq resource.ReadRequest, rs *resource.ReadResponse) {
	var data tagResourceData
	if diag := rq.State.Get(ctx, &data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	api := uptimeapi.ClientWithResponses{ClientInterface: t.api}
	obj, err := api.GetServiceTagDetailWithResponse(ctx, *data.ID)
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

func (t *tagResourceImpl) Update(ctx context.Context, rq resource.UpdateRequest, rs *resource.UpdateResponse) {
	data := new(tagResourceData)
	if diag := rq.State.Get(ctx, data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	id := *data.ID
	if diag := rq.Config.Get(ctx, data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	res, err := t.api.PutServiceTagDetail(ctx, id, data.ToAPI())
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
	data.FromAPI(obj.Result)

	if res.StatusCode != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", res.Status)
		return
	}
	if diag := rs.State.Set(ctx, data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	return
}

func (t *tagResourceImpl) Delete(ctx context.Context, rq resource.DeleteRequest, rs *resource.DeleteResponse) {
	var data tagResourceData
	if diag := rq.State.Get(ctx, &data); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	_, err := t.api.DeleteServiceTagDetail(ctx, *data.ID)
	if err != nil {
		rs.Diagnostics.AddError("Delete failed", err.Error())
		return
	}
}

func (t *tagResourceImpl) ImportState(ctx context.Context, rq resource.ImportStateRequest, rs *resource.ImportStateResponse) {
	api := uptimeapi.ClientWithResponses{ClientInterface: t.api}
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
