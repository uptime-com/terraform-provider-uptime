package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewTagResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[tagResourceModel, upapi.Tag, upapi.Tag]{
		api: &tagResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "tag",
			Schema:         tagResourceSchema,
		},
	}
}

var tagResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id":  IDAttribute(),
		"url": URLAttribute(),
		"tag": schema.StringAttribute{
			Required: true,
		},
		"color_hex": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	},
}

type tagResourceModel struct {
	ID       types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL      types.String `tfsdk:"url" ref:"URL,opt"`
	Tag      types.String `tfsdk:"tag"`
	ColorHex types.String `tfsdk:"color_hex"`
}

var _ genericResourceAPI[upapi.Tag, upapi.Tag] = (*tagResourceAPI)(nil)

type tagResourceAPI struct {
	provider *providerImpl
}

func (c *tagResourceAPI) Create(ctx context.Context, arg upapi.Tag) (*upapi.Tag, error) {
	obj, err := c.provider.api.Tags().Create(ctx, arg)
	return obj, err
}

func (c *tagResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Tag, error) {
	obj, err := c.provider.api.Tags().Get(ctx, pk)
	return obj, err
}

func (c *tagResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.Tag) (*upapi.Tag, error) {
	obj, err := c.provider.api.Tags().Update(ctx, pk, arg)
	return obj, err
}

func (c *tagResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Tags().Delete(ctx, pk)
}

//func (r *tagResourceImpl) ImportState(ctx context.Context, rq resource.ImportStateRequest, rs *resource.ImportStateResponse) {
//	panic("implement me")
//	//api := uptimeapi.ClientWithResponses{ClientInterface: r.api}
//	//obj, err := api.GetServicetaglistWithResponse(ctx, &uptimeapi.GetServicetaglistParams{
//	//	Search: &rq.ID,
//	//})
//	//if err != nil {
//	//	rs.Diagnostics.AddError("Import failed", err.Error())
//	//	return
//	//}
//	//if obj.StatusCode() != http.StatusOK {
//	//	rs.Diagnostics.AddError("Bad response status", obj.Status())
//	//	return
//	//}
//	//for _, tag := range *obj.JSON200.Results {
//	//	if tag.Tag != rq.ID {
//	//		continue
//	//	}
//	//	data := new(tagResourceData)
//	//	data.FromAPI(tag)
//	//	if diag := rs.State.Set(ctx, data); diag.HasError() {
//	//		rs.Diagnostics.Append(diag...)
//	//		return
//	//	}
//	//	return
//	//}
//	//rs.Diagnostics.AddError("Import failed", "tag not found")
//	//return
//}
