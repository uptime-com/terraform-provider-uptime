package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckBlacklistResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[CheckBlacklistResourceModel, upapi.CheckBlacklist, upapi.Check](
		CheckBlacklistResourceAPI{provider: p},
		CheckBlacklistResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_blacklist",
			Schema: schema.Schema{
				Description: "Checks your domain against approximately 100 of the most well-known spam blacklists once per day to see if it's included on those lists. Import using the check ID: `terraform import uptime_check_blacklist.example 123`",
				Attributes: map[string]schema.Attribute{
					"id":             IDSchemaAttribute(),
					"url":            URLSchemaAttribute(),
					"name":           NameSchemaAttribute(),
					"address":        AddressHostnameSchemaAttributeDescription("Domain name to check"),
					"contact_groups": ContactGroupsSchemaAttribute(),
					"locations":      LocationsReadOnlySchemaAttribute(),
					"tags":           TagsSchemaAttribute(),
					"is_paused":      IsPausedSchemaAttribute(),
					"num_retries":    NumRetriesAttribute(2),
					"notes":          NotesSchemaAttribute(),
				},
			},
		},
		ImportStateSimpleID,
	)
}

type CheckBlacklistResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
}

func (a CheckBlacklistResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckBlacklistResourceModel, diag.Diagnostics) {
	model := *new(CheckBlacklistResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckBlacklistResourceModelAdapter) ToAPIArgument(model CheckBlacklistResourceModel) (_ *upapi.CheckBlacklist, err error) {
	return &upapi.CheckBlacklist{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroups(model.ContactGroups),
		Locations:     a.Locations(model.Locations),
		Tags:          a.Tags(model.Tags),
		IsPaused:      model.IsPaused.ValueBool(),
		Address:       model.Address.ValueString(),
		NumRetries:    model.NumRetries.ValueInt64(),
		Notes:         model.Notes.ValueString(),
	}, nil
}

func (a CheckBlacklistResourceModelAdapter) FromAPIResult(api upapi.Check) (_ *CheckBlacklistResourceModel, err error) {
	model := CheckBlacklistResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsValue(api.ContactGroups),
		Locations:     a.LocationsValue(api.Locations),
		Tags:          a.TagsValue(api.Tags),
		IsPaused:      types.BoolValue(api.IsPaused),
		Address:       types.StringValue(api.Address),
		NumRetries:    types.Int64Value(api.NumRetries),
		Notes:         types.StringValue(api.Notes),
	}
	return &model, nil
}

type CheckBlacklistResourceModel struct {
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

func (m CheckBlacklistResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

var _ API[upapi.CheckBlacklist, upapi.Check] = (*CheckBlacklistResourceAPI)(nil)

type CheckBlacklistResourceAPI struct {
	provider *providerImpl
}

func (a CheckBlacklistResourceAPI) Create(ctx context.Context, arg upapi.CheckBlacklist) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateBlacklist(ctx, arg)
}

func (a CheckBlacklistResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a CheckBlacklistResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckBlacklist) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateBlacklist(ctx, pk, arg)
}

func (a CheckBlacklistResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
