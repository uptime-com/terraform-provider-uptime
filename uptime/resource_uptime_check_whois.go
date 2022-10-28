package uptime

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/uptime-com/uptime-client-go"
)

func resourceUptimeCheckWhois() *schema.Resource {
	var check WhoisCheck
	return &schema.Resource{
		Create: checkCreateFunc(check),
		Read:   checkReadFunc(check),
		Update: checkUpdateFunc(check),
		Delete: checkDeleteFunc(check),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Required attributes: Common
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDomain,
			},
			"contact_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Required attributes: Specific
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
			"expect_string": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func validateDomain(val interface{}, key string) (warns []string, errs []error) {
	urlStr := val.(string)

	if govalidator.IsDNSName(urlStr) != true {
		errs = append(errs, fmt.Errorf("Invalid domain: %s", urlStr))
	}

	return
}

// WhoisCheck implements the CheckType interface for Uptime.com Whois/Domain Expiry checks.
type WhoisCheck struct {
	CheckType
}

func (WhoisCheck) typeStr() string { return "WHOIS" }

func (WhoisCheck) getSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	if attr, ok := d.GetOk("days_before_expiry"); ok {
		c.Threshold = attr.(int)
	}

	if attr, ok := d.GetOk("expect_string"); ok {
		c.ExpectString = attr.(string)
	}
}

func (WhoisCheck) setSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	d.Set("days_before_expiry", c.Threshold)
	d.Set("expect_string", c.ExpectString)
}
