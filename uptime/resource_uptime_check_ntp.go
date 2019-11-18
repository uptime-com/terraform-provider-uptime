package uptime

import (
	"fmt"

	uptime "github.com/uptime-com/rest-api-clients/golang/uptime"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUptimeCheckNTP() *schema.Resource {
	return &schema.Resource{
		Create: checkCreateFunc(ntpCheck),
		Read: checkReadFunc(ntpCheck),
		Update: checkUpdateFunc(ntpCheck),
		Delete: checkDeleteFunc(ntpCheck),
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

			// Required attributes: Specific
			"interval": {
				Type: schema.TypeInt,
				Required: true,
			},
			"locations": {
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
			"include_in_global_metrics": {
				Type: schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			// Optional attributes: Specific
			"ip_version": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func (val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					valid := map[string]bool{
						"IPV4": true,
						"IPV6": true,
					}
					if _, ok := valid[v]; !ok {
						errs = append(errs, fmt.Errorf("Invalid IP version %s. Choose one of: IPV4, IPV6", v))
					}
					return
				},
			},
			"sensitivity": {
				Type: schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"threshold": {
				Type: schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type: schema.TypeInt,
				Optional: true,
				Computed: true,
			},

		},
	}
}

// NTPCheck implements the CheckType interface for Uptime.com NTP checks.
type NTPCheck struct{}

func (NTPCheck) typeStr() string {return "NTP"}

func (NTPCheck) getSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	if attr, ok := d.GetOk("interval"); ok {
		c.Interval = attr.(int)
	}
	if attr, ok := d.GetOk("ip_version"); ok {
		c.IPVersion = attr.(string)
	}
	if attr, ok := d.GetOk("locations"); ok {
		c.Locations = expandSetAttr(attr)
	}
	if attr, ok := d.GetOk("port"); ok {
		c.Port = attr.(int)
	}
	if attr, ok := d.GetOk("sensitivity"); ok {
		c.Sensitivity = attr.(int)
	}
	if attr, ok := d.GetOk("threshold"); ok {
		c.Threshold = attr.(int)
	}
}

func (NTPCheck) setSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	d.Set("interval", c.Interval)
	d.Set("locations", c.Locations)
	d.Set("ip_version", c.IPVersion)
	d.Set("sensitivity", c.Sensitivity)
	d.Set("threshold", c.Threshold)
	d.Set("port", c.Port)
}

var ntpCheck NTPCheck
