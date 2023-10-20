package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckWHOISResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkWHOISResourceModel, upapi.CheckWHOIS, upapi.Check]{
		api: &checkWHOISResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_whois",
			Schema:         checkWHOISResourceSchema,
		},
	}
}

var checkWHOISResourceSchema = schema.Schema{
	Description: "Monitor domain's expiry date and registration details",
	Attributes: map[string]schema.Attribute{
		"id":             IDAttribute(),
		"url":            URLAttribute(),
		"name":           NameAttribute(),
		"contact_groups": ContactGroupsAttribute(),
		"locations":      LocationsReadOnlyAttribute(),
		"tags":           TagsAttribute(),
		"is_paused":      IsPausedAttribute(),
		"threshold":      ThresholdDescriptionAttribute(20, "Raise an alert if there are less than this many days before the domain needs to be renewed."),
		"num_retries":    NumRetriesAttribute(2),
		"notes":          NotesAttribute(),

		"address": AddressHostnameAttribute(),
		"expect_string": schema.StringAttribute{
			Required:    true,
			Description: "The current domain registration info that should always match.",
		},
	},
}

type checkWHOISResourceModel struct {
	ID            types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL           types.String `tfsdk:"url" ref:"URL,opt"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	Locations     types.Set    `tfsdk:"locations"`
	Tags          types.Set    `tfsdk:"tags"`
	IsPaused      types.Bool   `tfsdk:"is_paused"`
	Address       types.String `tfsdk:"address"`
	ExpectString  types.String `tfsdk:"expect_string"`
	Threshold     types.Int64  `tfsdk:"threshold"`
	NumRetries    types.Int64  `tfsdk:"num_retries"`
	Notes         types.String `tfsdk:"notes"`
}

var _ genericResourceAPI[upapi.CheckWHOIS, upapi.Check] = (*checkWHOISResourceAPI)(nil)

type checkWHOISResourceAPI struct {
	provider *providerImpl
}

func (c *checkWHOISResourceAPI) Create(ctx context.Context, arg upapi.CheckWHOIS) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateWHOIS(ctx, arg)
}

func (c *checkWHOISResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c *checkWHOISResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckWHOIS) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateWHOIS(ctx, pk, arg)
}

func (c *checkWHOISResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
