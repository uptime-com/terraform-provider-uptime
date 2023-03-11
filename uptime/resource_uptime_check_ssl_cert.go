package uptime

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/uptime-com/uptime-client-go"
)

func resourceUptimeCheckSSLCert() *schema.Resource {
	return &schema.Resource{
		Create: checkCreateFunc(sslCheck),
		Read:   checkReadFunc(sslCheck),
		Update: checkUpdateFunc(sslCheck),
		Delete: checkDeleteFunc(sslCheck),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Required attributes: Common
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Required attributes: Specific
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					p := val.(string)
					valid := map[string]bool{
						"https": true,
						"http":  true,
						"pop3":  true,
						"imap":  true,
						"ftp":   true,
						"xmpp":  true,
						"irc":   true,
						"ldap":  true,
					}
					if _, ok := valid[p]; !ok {
						errs = append(errs, fmt.Errorf("invalid protocol for SSL Cert check: %s", p))
					}
					return
				},
			},
			"days_before_expiry": {
				Type:     schema.TypeInt,
				Required: true,
			},

			// Optional attributes: Common
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},

			// Optional attributes: Specific
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// SSLCheck implements the CheckType interface for Uptime.com SSL Cert checks.
type SSLCheck struct{}

func (SSLCheck) typeStr() string { return "SSL_CERT" }

func (SSLCheck) getSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	if attr, ok := d.GetOk("protocol"); ok {
		c.Protocol = attr.(string)
	}

	if attr, ok := d.GetOk("days_before_expiry"); ok {
		c.Threshold = attr.(int)
	}

	if attr, ok := d.GetOk("port"); ok {
		c.Port = attr.(int)
	}
}

func (SSLCheck) setSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	d.Set("protocol", c.Protocol)
	d.Set("days_before_expiry", c.Threshold)
	d.Set("port", c.Port)
}

var sslCheck SSLCheck
