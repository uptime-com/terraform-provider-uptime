package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckBlacklistResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkBlacklistResourceModel, upapi.CheckBlacklist, upapi.Check]{
		api: &checkBlacklistResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_blacklist",
			Schema:         checkBlacklistResourceSchema,
		},
	}
}

var checkBlacklistResourceSchema = schema.Schema{
	Description: "Checks your domain against approximately 100 of the most well-known spam blacklists once per day to see if it’s included on those lists",
	Attributes: map[string]schema.Attribute{
		"id":             IDAttribute(),
		"url":            URLAttribute(),
		"name":           NameAttribute(),
		"contact_groups": ContactGroupsAttribute(),
		"locations":      LocationsReadOnlyAttribute(),
		"tags":           TagsAttribute(),
		"is_paused":      IsPausedAttribute(),
		"num_retries":    NumRetriesAttribute(2),
		"notes":          NotesAttribute(),

		"address": schema.StringAttribute{
			Required:    true,
			Description: "Domain name to check",
		},
	},
}

type checkBlacklistResourceModel struct {
	ID            types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL           types.String `tfsdk:"url" ref:"URL,opt"`
	Name          types.String `tfsdk:"name"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	Locations     types.Set    `tfsdk:"locations"`
	Tags          types.Set    `tfsdk:"tags"`
	IsPaused      types.Bool   `tfsdk:"is_paused"`
	Address       types.String `tfsdk:"address"`
	NumRetries    types.Int64  `tfsdk:"num_retries"`
	Notes         types.String `tfsdk:"notes"`
}

var _ genericResourceAPI[upapi.CheckBlacklist, upapi.Check] = (*checkBlacklistResourceAPI)(nil)

type checkBlacklistResourceAPI struct {
	provider *providerImpl
}

func (a *checkBlacklistResourceAPI) Create(ctx context.Context, arg upapi.CheckBlacklist) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateBlacklist(ctx, arg)
}

func (a *checkBlacklistResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a *checkBlacklistResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckBlacklist) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateBlacklist(ctx, pk, arg)
}

func (a *checkBlacklistResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
