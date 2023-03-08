package provider

import (
	"bytes"
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/terraform-provider-uptime/internal/uptimeapi"
)

type checkAPIResourceData struct {
	ID                     types.Int64  `tfsdk:"id"                        map:"Pk"`
	Script                 types.String `tfsdk:"script"                    map:"MspScript"`
	Interval               types.Int64  `tfsdk:"interval"                  map:"MspInterval"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"            map:"ContactGroups"`
	Locations              types.Set    `tfsdk:"locations"                 map:"Locations"`
	Tags                   types.Set    `tfsdk:"tags"                      map:"Tags"`
	Name                   types.String `tfsdk:"name"                      map:"Name"`
	Threshold              types.Int64  `tfsdk:"threshold"                 map:"MspThreshold"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"               map:"MspSensitivity"`
	Notes                  types.String `tfsdk:"notes"                     map:"MspNotes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics" map:"MspIncludeInGlobalMetrics"`
}

func (d *checkAPIResourceData) ToChecks(ctx context.Context) (obj uptimeapi.Checks, diags diag.Diagnostics) {
	diags = fromTerraform(ctx, &obj, d)
	return
}

func (d *checkAPIResourceData) ToChecksAPI(ctx context.Context) (obj uptimeapi.ChecksAPI, diags diag.Diagnostics) {
	diags = fromTerraform(ctx, &obj, d)
	return
}

func (d *checkAPIResourceData) From(obj any) diag.Diagnostics {
	return toTerraform(d, obj)
}

type checkAPIResourceImpl struct {
	api uptimeapi.ClientInterface
}

func (r *checkAPIResourceImpl) apiWithResponses() uptimeapi.ClientWithResponsesInterface {
	return &uptimeapi.ClientWithResponses{ClientInterface: r.api}
}

func (r *checkAPIResourceImpl) Metadata(_ context.Context, rq resource.MetadataRequest, rs *resource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_check_api"
}

func (r *checkAPIResourceImpl) Schema(_ context.Context, rq resource.SchemaRequest, rs *resource.SchemaResponse) {
	rs.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"contact_groups": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"interval": schema.Int64Attribute{
				Required: true,
			},
			"locations": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"script": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"notes": schema.StringAttribute{
				Optional: true,
			},
			"include_in_global_metrics": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"sensitivity": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"threshold": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *checkAPIResourceImpl) Create(ctx context.Context, rq resource.CreateRequest, rs *resource.CreateResponse) {
	var tfData checkAPIResourceData
	if diags := rq.Plan.Get(ctx, &tfData); diags.HasError() {
		rs.Diagnostics = diags
		return
	}
	apiData, diags := tfData.ToChecksAPI(ctx)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	res, err := r.apiWithResponses().PostServiceCreateApiWithResponse(ctx, apiData)
	if err != nil {
		rs.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	if res.StatusCode() != 200 {
		rs.Diagnostics.AddError("Bad response status", prettyResponse(res.Status(), bytes.NewReader(res.Body)))
		return
	}
	diags = tfData.From(*res.JSON200)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	diags = rs.State.Set(ctx, tfData)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	return
}

func (r *checkAPIResourceImpl) Read(ctx context.Context, rq resource.ReadRequest, rs *resource.ReadResponse) {
	data := checkAPIResourceData{}
	diags := rq.State.Get(ctx, &data)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	res, err := r.apiWithResponses().GetServiceDetailWithResponse(ctx, data.ID.String())
	if err != nil {
		rs.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	if res.StatusCode() != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", prettyResponse(res.Status(), bytes.NewReader(res.Body)))
		return
	}
	diags = data.From(*res.JSON200)
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

func (r *checkAPIResourceImpl) Update(ctx context.Context, rq resource.UpdateRequest, rs *resource.UpdateResponse) {
	prev := new(checkAPIResourceData)
	if diag := rq.State.Get(ctx, prev); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	next := new(checkAPIResourceData)
	if diag := rq.Config.Get(ctx, next); diag.HasError() {
		rs.Diagnostics.Append(diag...)
		return
	}
	obj, diags := next.ToChecks(ctx)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	res, err := r.apiWithResponses().PutServiceDetailWithResponse(ctx, prev.ID.String(), obj)
	if err != nil {
		rs.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	if res.StatusCode() != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", prettyResponse(res.Status(), bytes.NewReader(res.Body)))
		return
	}
	data := checkAPIResourceData{}
	diags = data.From(*res.JSON200)
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

func (r *checkAPIResourceImpl) Delete(ctx context.Context, rq resource.DeleteRequest, rs *resource.DeleteResponse) {
	data := checkAPIResourceData{}
	diags := rq.State.Get(ctx, &data)
	if diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	res, err := r.api.DeleteServiceDetail(ctx, data.ID.String())
	if err != nil {
		rs.Diagnostics.AddError("Delete failed", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		rs.Diagnostics.AddError("Bad response status", prettyResponse(res.Status, res.Body))
		return
	}
	return
}
