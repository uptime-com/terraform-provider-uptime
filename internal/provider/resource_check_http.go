package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckHTTPResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkHTTPResourceModel, upapi.CheckHTTP, upapi.Check]{
		api: &checkHTTPResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_http",
			Schema:         checkHTTPResourceSchema,
		},
	}
}

var checkHTTPResourceSchema = schema.Schema{
	Description: "Monitor a URL for specific status code(s)",
	Attributes: map[string]schema.Attribute{
		"id":                        IDAttribute(),
		"url":                       URLAttribute(),
		"name":                      NameAttribute(),
		"contact_groups":            ContactGroupsAttribute(),
		"locations":                 LocationsAttribute(),
		"tags":                      TagsAttribute(),
		"is_paused":                 IsPausedAttribute(),
		"interval":                  IntervalAttribute(5),
		"threshold":                 ThresholdAttribute(40),
		"sensitivity":               SensitivityAttribute(2),
		"num_retries":               NumRetriesAttribute(2),
		"notes":                     NotesAttribute(),
		"include_in_global_metrics": IncludeInGlobalMetricsAttribute(),

		"address": schema.StringAttribute{
			Required:    true,
			Description: "URL to check",
		},
		"port": schema.Int64Attribute{
			Computed: true,
			Optional: true,
			Default:  int64default.StaticInt64(0),
		},
		"username": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(""),
		},
		"password": schema.StringAttribute{
			Optional:  true,
			Computed:  true,
			Sensitive: true,
			Default:   stringdefault.StaticString(""),
		},
		"proxy": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(""),
		},
		"status_code": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("200"),
		},
		"send_string": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			Description: "String to post",
		},
		"expect_string": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(""),
		},
		"expect_string_type": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("STRING"),
			Validators: []validator.String{
				OneOfStringValidator([]string{"STRING", "REGEX", "INVERSE_REGEX"}),
			},
		},
		"headers": schema.MapAttribute{
			ElementType: types.ListType{
				ElemType: types.StringType,
			},
			Optional: true,
			Computed: true,
			Default:  mapdefault.StaticValue(types.MapValueMust(types.ListType{ElemType: types.StringType}, map[string]attr.Value{})),
		},
		"version": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(2),
			Description: "Check version to use. Keep default value unless you are absolutely sure you need to change it",
		},

		"encryption": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("SSL_TLS"),
			Description: "Whether to verify SSL/TLS certificates",
		},
	},
}

type checkHTTPResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.Set    `tfsdk:"contact_groups"`
	Locations              types.Set    `tfsdk:"locations"`
	Tags                   types.Set    `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Interval               types.Int64  `tfsdk:"interval"`
	Address                types.String `tfsdk:"address"`
	Port                   types.Int64  `tfsdk:"port"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	Proxy                  types.String `tfsdk:"proxy"`
	StatusCode             types.String `tfsdk:"status_code"`
	SendString             types.String `tfsdk:"send_string"`
	ExpectString           types.String `tfsdk:"expect_string"`
	ExpectStringType       types.String `tfsdk:"expect_string_type"`
	Encryption             types.String `tfsdk:"encryption"`
	Threshold              types.Int64  `tfsdk:"threshold"`
	Headers                types.Map    `tfsdk:"headers" ref:",extra=headers"`
	Version                types.Int64  `tfsdk:"version"`
	Sensitivity            types.Int64  `tfsdk:"sensitivity"`
	NumRetries             types.Int64  `tfsdk:"num_retries"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
}

var _ genericResourceAPI[upapi.CheckHTTP, upapi.Check] = (*checkHTTPResourceAPI)(nil)

type checkHTTPResourceAPI struct {
	provider *providerImpl
}

func (a *checkHTTPResourceAPI) Create(ctx context.Context, arg upapi.CheckHTTP) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateHTTP(ctx, arg)
}

func (a *checkHTTPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a *checkHTTPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckHTTP) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateHTTP(ctx, pk, arg)
}

func (a *checkHTTPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
