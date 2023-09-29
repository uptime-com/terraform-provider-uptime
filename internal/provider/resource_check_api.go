package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckAPIResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkAPIResourceModel, upapi.CheckAPI, upapi.Check]{
		api: &checkAPIResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_api",
			Schema:         checkAPIResourceSchema,
		},
	}
}

var checkAPIResourceSchema = schema.Schema{
	Description: "Multi-step advanced check type that is intended to monitor API such as REST or SOAP",
	Attributes: map[string]schema.Attribute{
		"id":                        IDAttribute(),
		"url":                       URLAttribute(),
		"name":                      NameAttribute(),
		"contact_groups":            ContactGroupsAttribute(),
		"locations":                 LocationsAttribute(),
		"tags":                      TagsAttribute(),
		"is_paused":                 IsPausedAttribute(),
		"interval":                  IntervalAttribute(5),
		"threshold":                 ThresholdAttribute(30),
		"sensitivity":               SensitivityAttribute(2),
		"num_retries":               NumRetriesAttribute(2),
		"notes":                     NotesAttribute(),
		"include_in_global_metrics": IncludeInGlobalMetricsAttribute(),
		"response_time_sla":         ResponseTimeSLAAttribute("30s"),

		"script": schema.StringAttribute{
			Required: true,
		},
	},
}

type checkAPIResourceModel struct {
	ID                     types.Int64          `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String         `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String         `tfsdk:"name"`
	ContactGroups          types.Set            `tfsdk:"contact_groups"`
	Locations              types.Set            `tfsdk:"locations"`
	Tags                   types.Set            `tfsdk:"tags"`
	IsPaused               types.Bool           `tfsdk:"is_paused"`
	Interval               types.Int64          `tfsdk:"interval"`
	Threshold              types.Int64          `tfsdk:"threshold"`
	Script                 types.String         `tfsdk:"script"`
	Sensitivity            types.Int64          `tfsdk:"sensitivity"`
	NumRetries             types.Int64          `tfsdk:"num_retries"`
	Notes                  types.String         `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool           `tfsdk:"include_in_global_metrics"`
	ResponseTimeSLA        customtypes.Duration `tfsdk:"response_time_sla"`
}

var _ genericResourceAPI[upapi.CheckAPI, upapi.Check] = (*checkAPIResourceAPI)(nil)

type checkAPIResourceAPI struct {
	provider *providerImpl
}

func (a *checkAPIResourceAPI) Create(ctx context.Context, arg upapi.CheckAPI) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateAPI(ctx, arg)
}

func (a *checkAPIResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a *checkAPIResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckAPI) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateAPI(ctx, pk, arg)
}

func (a *checkAPIResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
