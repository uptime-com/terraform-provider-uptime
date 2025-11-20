package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewSubaccountResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[SubaccountResourceModel, upapi.SubaccountCreateRequest, upapi.Subaccount]{
		api: &SubaccountResourceAPI{provider: p},
		mod: SubaccountResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "subaccount",
			Schema: schema.Schema{
				Description: "Manage Uptime.com subaccounts.\n\n" +
					"**IMPORTANT:** This resource requires the subaccounts feature to be enabled for your account. " +
					"Attempts to create subaccounts without this feature enabled will fail with a PERMISSION_DENIED error.",
				Attributes: map[string]schema.Attribute{
					"id":   IDSchemaAttribute(),
					"url":  URLSchemaAttribute(),
					"name": NameSchemaAttribute(),
				},
			},
		},
	}
}

type SubaccountResourceModel struct {
	ID   types.Int64  `tfsdk:"id"`
	URL  types.String `tfsdk:"url"`
	Name types.String `tfsdk:"name"`
}

func (m SubaccountResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type SubaccountResourceModelAdapter struct{}

func (a SubaccountResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*SubaccountResourceModel, diag.Diagnostics) {
	model := *new(SubaccountResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a SubaccountResourceModelAdapter) ToAPIArgument(model SubaccountResourceModel) (*upapi.SubaccountCreateRequest, error) {
	return &upapi.SubaccountCreateRequest{
		Name: model.Name.ValueString(),
	}, nil
}

func (a SubaccountResourceModelAdapter) FromAPIResult(api upapi.Subaccount) (*SubaccountResourceModel, error) {
	return &SubaccountResourceModel{
		ID:   types.Int64Value(api.PK),
		URL:  types.StringValue(api.URL),
		Name: types.StringValue(api.Name),
	}, nil
}

type SubaccountResourceAPI struct {
	provider *providerImpl
}

func (a SubaccountResourceAPI) Create(ctx context.Context, arg upapi.SubaccountCreateRequest) (*upapi.Subaccount, error) {
	return a.provider.api.Subaccounts().Create(ctx, arg)
}

func (a SubaccountResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Subaccount, error) {
	return a.provider.api.Subaccounts().Get(ctx, pk)
}

func (a SubaccountResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.SubaccountCreateRequest) (*upapi.Subaccount, error) {
	// For update, we need to convert CreateRequest to UpdateRequest
	return a.provider.api.Subaccounts().Update(ctx, pk, upapi.SubaccountUpdateRequest{
		Name: arg.Name,
	})
}

func (a SubaccountResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	// Note: The Subaccounts endpoint doesn't have a Delete method in the API
	// This is a limitation of the current API - subaccounts cannot be deleted via API
	// Return an error to inform users
	return fmt.Errorf("subaccounts cannot be deleted via the API - this is a limitation of the Uptime.com API. Please manually delete the subaccount via the web interface if needed")
}
