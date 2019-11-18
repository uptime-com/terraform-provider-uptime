package uptime

import (
	"github.com/hashicorp/terraform/helper/schema"
	uptime "github.com/uptime-com/rest-api-clients/golang/uptime"
)

func resourceUptimeCheckDomainBlacklist() *schema.Resource {
	return &schema.Resource{

		Create: checkCreateFunc(blacklistCheck),
		Read: checkReadFunc(blacklistCheck),
		Update: checkUpdateFunc(blacklistCheck),
		Delete: checkDeleteFunc(blacklistCheck),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Required attributes: Common
			"address": {
				Type: schema.TypeString,
				Required: true,
			},
			"contact_groups": {
				Type: schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Optional attributes: Common
			"name": {
				Type: schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type: schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notes": {
				Type: schema.TypeString,
				Optional: true,
			Default: "Managed by Terraform",
			},
		},
	}
}

// BlacklistCheck implements the CheckType interface for Uptime.com Domain Blacklist checks.
type BlacklistCheck struct{}

func (BlacklistCheck) typeStr() string {return "BLACKLIST"}

func (BlacklistCheck) getSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {}

func (BlacklistCheck) setSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {}

var blacklistCheck BlacklistCheck
