package uptime

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"uptime_tag": resourceUptimeCheckTag(),
			"uptime_check_dns": resourceUptimeCheckDNS(),
			"uptime_check_domain_blacklist": resourceUptimeCheckDomainBlacklist(),
			"uptime_check_http": resourceUptimeCheckHTTP(),
			"uptime_check_malware": resourceUptimeCheckMalware(),
			"uptime_check_ntp": resourceUptimeCheckNTP(),
			"uptime_check_ssl_cert": resourceUptimeCheckSSLCert(),
			"uptime_check_whois": resourceUptimeCheckWhois(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(data *schema.ResourceData) (interface{}, error) {
	c := Config{
		Token: data.Get("token").(string),
		RateMilliseconds: 500,
	}

	log.Println("[INFO] Initializing Uptime client")
	return c.Client()
}
