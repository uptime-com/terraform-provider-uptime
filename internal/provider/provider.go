package provider

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

var _ provider.Provider = (*providerImpl)(nil)

type providerImpl struct {
	api       upapi.API
	version   string
	locations []string
}

type providerConfig struct {
	Endpoint  types.String  `tfsdk:"endpoint"`
	Token     types.String  `tfsdk:"token"`
	RateLimit types.Float64 `tfsdk:"rate_limit"`
	Trace     types.Bool    `tfsdk:"trace"`
}

func (p *providerImpl) Metadata(_ context.Context, _ provider.MetadataRequest, rs *provider.MetadataResponse) {
	rs.TypeName = "uptime"
	rs.Version = p.version
}

func (p *providerImpl) Schema(_ context.Context, _ provider.SchemaRequest, rs *provider.SchemaResponse) {
	rs.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"rate_limit": schema.Float64Attribute{
				Optional:    true,
				Description: "The rate limit to use for API calls in requests per second, defaults to 0.5",
			},
			"trace": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (p *providerImpl) UserAgentString() string {
	return fmt.Sprintf("Uptime.com Terraform Provider %s %s/%s", p.version, runtime.GOOS, runtime.GOARCH)
}

func (p *providerImpl) Configure(ctx context.Context, rq provider.ConfigureRequest, rs *provider.ConfigureResponse) {
	var cfg providerConfig
	if diags := rq.Config.Get(ctx, &cfg); diags.HasError() {
		rs.Diagnostics.Append(diags...)
		return
	}
	if cfg.Token.IsNull() {
		cfg.Token = types.StringValue(os.Getenv("UPTIME_TOKEN"))
	}
	if cfg.Trace.IsNull() {
		cfg.Trace = types.BoolValue(os.Getenv("UPTIME_TRACE") != "")
	}
	if cfg.RateLimit.IsNull() {
		cfg.RateLimit = types.Float64Value(0.5)
	}
	opts := []upapi.Option{
		upapi.WithToken(cfg.Token.ValueString()),
		upapi.WithUserAgent(p.UserAgentString()),
		upapi.WithRateLimit(cfg.RateLimit.ValueFloat64()),
		upapi.WithRetry(10, time.Second*30, os.Stderr),
	}
	if ep := cfg.Endpoint.ValueString(); ep != "" {
		opts = append(opts, upapi.WithBaseURL(ep))
	}
	if cfg.Trace.ValueBool() {
		opts = append(opts, upapi.WithTrace(os.Stderr))
	}
	api, err := upapi.New(opts...)
	if err != nil {
		rs.Diagnostics.AddError("Failed to initialize API client", err.Error())
		return
	}
	p.api = api
	servers, err := p.api.ProbeServers().List(ctx)
	if err != nil {
		rs.Diagnostics.AddError("Failed to get list of locations", err.Error())
		return

	}
	locationsSet := make(map[string]struct{}, len(servers))
	for _, server := range servers {
		if _, ok := locationsSet[server.Location]; !ok {
			locationsSet[server.Location] = struct{}{}
			p.locations = append(p.locations, server.Location)
		}
	}
}

func (p *providerImpl) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return NewLocationsDataSource(ctx, p) },
	}
}

func (p *providerImpl) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return NewCheckAPIResource(ctx, p) },
		func() resource.Resource { return NewCheckBlacklistResource(ctx, p) },
		func() resource.Resource { return NewCheckDNSResource(ctx, p) },
		func() resource.Resource { return NewCheckHeartbeatResource(ctx, p) },
		func() resource.Resource { return NewCheckHTTPResource(ctx, p) },
		func() resource.Resource { return NewCheckICMPResource(ctx, p) },
		func() resource.Resource { return NewCheckMalwareResource(ctx, p) },
		func() resource.Resource { return NewCheckNTPResource(ctx, p) },
		func() resource.Resource { return NewCheckSSLCertResource(ctx, p) },
		func() resource.Resource { return NewCheckWHOISResource(ctx, p) },

		func() resource.Resource { return NewContactResource(ctx, p) },
		func() resource.Resource { return NewStatusPageResource(ctx, p) },
		func() resource.Resource { return NewTagResource(ctx, p) },
	}
}

func VersionFactory(version string) func() provider.Provider {
	return func() provider.Provider {
		return &providerImpl{
			version: version,
		}
	}
}
