package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckDNSResource(_ context.Context, p *providerImpl) resource.Resource {
	return &genericResource[checkDNSResourceModel, upapi.CheckDNS, upapi.Check]{
		api: &checkDNSResourceAPI{provider: p},
		metadata: genericResourceMetadata{
			TypeNameSuffix: "check_dns",
			Schema:         checkDNSResourceSchema,
		},
	}
}

var checkDNSResourceSchema = schema.Schema{
	Description: "Monitor for DNS failures or changes",
	Attributes: map[string]schema.Attribute{
		"id":             IDAttribute(),
		"url":            URLAttribute(),
		"name":           NameAttribute(),
		"contact_groups": ContactGroupsAttribute(),
		"locations":      LocationsAttribute(),
		"tags":           TagsAttribute(),
		"is_paused":      IsPausedAttribute(),
		"interval":       IntervalAttribute(5),
		"threshold":      ThresholdAttribute(20),
		"address":        AddressHostnameAttribute(),
		"dns_server": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			Description: "DNS server to query",
		},
		"dns_record_type": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("ANY"),
			Validators: []validator.String{
				OneOfStringValidator([]string{"ANY", "A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT"}),
			},
		},
		"expect_string": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
			Description: "IP Address, Domain Name or String to expect in response",
		},
		"sensitivity":               SensitivityAttribute(2),
		"num_retries":               NumRetriesAttribute(2),
		"notes":                     NotesAttribute(),
		"include_in_global_metrics": IncludeInGlobalMetricsAttribute(),
		"response_time_sla":         ResponseTimeSLAAttribute("1s"),
	},
}

type checkDNSResourceModel struct {
	ID                     types.Int64          `tfsdk:"id"  ref:"PK,opt"`
	URL                    types.String         `tfsdk:"url" ref:"URL,opt"`
	Name                   types.String         `tfsdk:"name"`
	ContactGroups          types.Set            `tfsdk:"contact_groups"`
	Locations              types.Set            `tfsdk:"locations"`
	Tags                   types.Set            `tfsdk:"tags"`
	IsPaused               types.Bool           `tfsdk:"is_paused"`
	Interval               types.Int64          `tfsdk:"interval"`
	Threshold              types.Int64          `tfsdk:"threshold"`
	Address                types.String         `tfsdk:"address"`
	DNSServer              types.String         `tfsdk:"dns_server"`
	DNSRecordType          types.String         `tfsdk:"dns_record_type"`
	ExpectString           types.String         `tfsdk:"expect_string"`
	Sensitivity            types.Int64          `tfsdk:"sensitivity"`
	NumRetries             types.Int64          `tfsdk:"num_retries"`
	Notes                  types.String         `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool           `tfsdk:"include_in_global_metrics"`
	ResponseTimeSLA        customtypes.Duration `tfsdk:"response_time_sla"`
}

var _ genericResourceAPI[upapi.CheckDNS, upapi.Check] = (*checkDNSResourceAPI)(nil)

type checkDNSResourceAPI struct {
	provider *providerImpl
}

func (a *checkDNSResourceAPI) Create(ctx context.Context, arg upapi.CheckDNS) (*upapi.Check, error) {
	return a.provider.api.Checks().CreateDNS(ctx, arg)
}

func (a *checkDNSResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return a.provider.api.Checks().Get(ctx, pk)
}

func (a *checkDNSResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckDNS) (*upapi.Check, error) {
	return a.provider.api.Checks().UpdateDNS(ctx, pk, arg)
}

func (a *checkDNSResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return a.provider.api.Checks().Delete(ctx, pk)
}
