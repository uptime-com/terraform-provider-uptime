package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uptime-com/terraform-provider-uptime/internal/uptimeapi"
)

type Option func(*providerImpl)

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

var _ provider.Provider = &providerImpl{}

type providerImpl struct {
	client  uptimeapi.ClientInterface
	version string
}

type cfgModel struct {
	Endpoint    types.String `tfsdk:"endpoint"`
	Token       types.String `tfsdk:"token"`
	RateLimitMs types.Int64  `tfsdk:"rate_limit_ms"`
}

type envModel struct {
	Endpoint    string `env:"UPTIME_ENDPOINT" envDefault:"https://uptime.com"`
	Token       string `env:"UPTIME_TOKEN"`
	RateLimitMs int64  `env:"UPTIME_RATE_LIMIT_MS" envDefault:"500"`
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
			"rate_limit_ms": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

func (p *providerImpl) Configure(ctx context.Context, rq provider.ConfigureRequest, rs *provider.ConfigureResponse) {
	envObj := envModel{}
	if err := env.Parse(&envObj); err != nil {
		rs.Diagnostics.AddError("Failed to parse environment variables", err.Error())
		return
	}
	cfgObj := cfgModel{}
	rs.Diagnostics.Append(rq.Config.Get(ctx, &cfgObj)...)
	if rs.Diagnostics.HasError() {
		return
	}

	var opts []uptimeapi.ClientOption

	var endpoint string
	if s := cfgObj.Endpoint.ValueString(); s != "" {
		endpoint = s
	} else {
		endpoint = envObj.Endpoint
	}

	var token string
	if s := cfgObj.Token.ValueString(); s != "" {
		token = s
	} else {
		token = envObj.Token
	}
	if token != "" {
		opts = append(opts, uptimeapi.WithToken(token))
	}

	var every time.Duration
	if i64 := cfgObj.RateLimitMs.ValueInt64(); i64 > 0 {
		every = time.Duration(i64) * time.Millisecond
	} else {
		every = time.Duration(envObj.RateLimitMs) * time.Millisecond
	}
	if every > 0 {
		opts = append(opts, uptimeapi.WithRateLimitEvery(every))
	}

	c, err := uptimeapi.NewClientWithResponses(endpoint, opts...)
	if err != nil {
		rs.Diagnostics.AddError("Failed to initialize API client", err.Error())
		return
	}
	p.client = c
}

func (p *providerImpl) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *providerImpl) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &tagResourceImpl{api: p.client} },
	}
}

func VersionFactory(version string) func() provider.Provider {
	return func() provider.Provider {
		p := providerImpl{}
		p.version = version
		return &p
	}
}
