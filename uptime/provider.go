package uptime

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider represents a resource provider in Terraform
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UPTIME_TOKEN", nil),
			},
			"rate_limit_ms": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UPTIME_RATE_LIMIT_MS", 500),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"uptime_tag":                    resourceUptimeCheckTag(),
			"uptime_check_dns":              resourceUptimeCheckDNS(),
			"uptime_check_domain_blacklist": resourceUptimeCheckDomainBlacklist(),
			"uptime_check_http":             resourceUptimeCheckHTTP(),
			"uptime_check_malware":          resourceUptimeCheckMalware(),
			"uptime_check_ntp":              resourceUptimeCheckNTP(),
			"uptime_check_ssl_cert":         resourceUptimeCheckSSLCert(),
			"uptime_check_whois":            resourceUptimeCheckWhois(),
			"uptime_check_heartbeat":        resourceUptimeCheckHeartbeat(),
			"uptime_integration_opsgenie":   resourceUptimeIntegrationOpsgenie(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c := Config{
		Token:            data.Get("token").(string),
		RateMilliseconds: data.Get("rate_limit_ms").(int),
	}

	log.Println("[INFO] Initializing Uptime client")

	cli, err := c.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return cli, nil
}
