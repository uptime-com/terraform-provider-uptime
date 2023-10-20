package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckICMPResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkICMPResourceModel, upapi.CheckICMP, upapi.Check]{
		api: &checkICMPResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_icmp",
			Schema:         checkICMPResourceSchema,
		},
	}
}

var checkICMPResourceSchema = schema.Schema{
	Description: "Monitor network activity for a specific domain or IP address",
	Attributes: map[string]schema.Attribute{
		"id":                        IDAttribute(),
		"url":                       URLAttribute(),
		"name":                      NameAttribute(),
		"contact_groups":            ContactGroupsAttribute(),
		"locations":                 LocationsAttribute(),
		"tags":                      TagsAttribute(),
		"is_paused":                 IsPausedAttribute(),
		"interval":                  IntervalAttribute(5),
		"num_retries":               NumRetriesAttribute(2),
		"use_ip_version":            UseIPVersionAttribute(),
		"notes":                     NotesAttribute(),
		"include_in_global_metrics": IncludeInGlobalMetricsAttribute(),
		"response_time_sla":         ResponseTimeSLAAttribute("1s"),

		"address": AddressHostnameAttribute(),
	},
}

type checkICMPResourceModel struct {
	ID                     types.Int64          `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String         `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String         `tfsdk:"name"`
	ContactGroups          types.Set            `tfsdk:"contact_groups"`
	Locations              types.Set            `tfsdk:"locations"`
	Tags                   types.Set            `tfsdk:"tags"`
	IsPaused               types.Bool           `tfsdk:"is_paused"`
	Interval               types.Int64          `tfsdk:"interval"`
	NumRetries             types.Int64          `tfsdk:"num_retries"`
	UseIPVersion           types.String         `tfsdk:"use_ip_version"`
	Notes                  types.String         `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool           `tfsdk:"include_in_global_metrics"`
	ResponseTimeSLA        customtypes.Duration `tfsdk:"response_time_sla"`

	Address types.String `tfsdk:"address"`
}

var _ genericResourceAPI[upapi.CheckICMP, upapi.Check] = (*checkICMPResourceAPI)(nil)

type checkICMPResourceAPI struct {
	provider *providerImpl
}

func (c *checkICMPResourceAPI) Create(ctx context.Context, arg upapi.CheckICMP) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateICMP(ctx, arg)
}

func (c *checkICMPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c *checkICMPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckICMP) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateICMP(ctx, pk, arg)
}

func (c *checkICMPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
