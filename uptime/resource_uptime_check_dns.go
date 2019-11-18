package uptime

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	uptime "github.com/uptime-com/rest-api-clients/golang/uptime"
)

func resourceUptimeCheckDNS() *schema.Resource {
	return &schema.Resource{
		Create: checkCreateFunc(dnsCheck),
		Read: checkReadFunc(dnsCheck),
		Update: checkUpdateFunc(dnsCheck),
		Delete: checkDeleteFunc(dnsCheck),
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
			"dns_record_type": {
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					valid := map[string]bool{
						"A": true,
						"AAAA": true,
						"CNAME": true,
						"MX": true,
						"NS": true,
						"PTR": true,
						"SOA": true,
						"TXT": true,
						"ANY": true,
					}
					if _, ok := valid[v]; !ok {
						errs = append(errs, fmt.Errorf("Invalid DNS Record Type %v", v))
					}
					return
				},
			},
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
			"dns_server": {
				Type: schema.TypeString,
				Optional: true,
			},
			"expect_string": {
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}
}

// DNSCheck implements the CheckType interface for Uptime.com DNS server checks.
type DNSCheck struct{}

func (DNSCheck) typeStr() string {return "DNS"}

func (DNSCheck) getSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	if attr, ok := d.GetOk("dns_record_type"); ok {
		c.DNSRecordType = attr.(string)
	}

	if attr, ok := d.GetOk("interval"); ok {
		c.Interval = attr.(int)
	}

	if attr, ok := d.GetOk("locations"); ok{
		c.Locations = expandSetAttr(attr)
	}

	if attr, ok := d.GetOk("sensitivity"); ok {
		c.Sensitivity = attr.(int)
	}

	if attr, ok := d.GetOk("threshold"); ok {
		c.Threshold = attr.(int)
	}

	if attr, ok := d.GetOk("dns_server"); ok {
		c.DNSServer = attr.(string)
	}

	if attr, ok := d.GetOk("expect_string"); ok {
		c.ExpectString = attr.(string)
	}
}

func (DNSCheck) setSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	d.Set("dns_record_type", c.DNSRecordType)
	d.Set("interval", c.Interval)
	d.Set("locations", c.Locations)
	d.Set("sensitivity", c.Sensitivity)
	d.Set("threshold", c.Threshold)
	d.Set("dns_server", c.DNSServer)
	d.Set("expect_string", c.ExpectString)
}

var dnsCheck DNSCheck
