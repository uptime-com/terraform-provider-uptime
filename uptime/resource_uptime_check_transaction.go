package uptime

import (
	uptime "github.com/uptime-com/rest-api-clients/golang/uptime"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUptimeCheckTransaction() *schema.Resource {
	return &schema.Resource{
		Create: checkCreateFunc(transactionCheck),
		Read: checkReadFunc(transactionCheck),
		Update: checkUpdateFunc(transactionCheck),
		Delete: checkDeleteFunc(transactionCheck),
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
			"script": {
				Type: schema.TypeString,
				Required: true,
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
		},
	}
}

// TransactionCheck implements the CheckType interface for Uptime.com Transaction checks.
type TransactionCheck struct{}

func (TransactionCheck) typeStr() string {return "TRANSACTION"}

func (TransactionCheck) getSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	if attr, ok := d.GetOk("interval"); ok {
		c.Interval = attr.(int)
	}
	if attr, ok := d.GetOk("locations"); ok {
		c.Locations = expandSetAttr(attr)
	}
	if attr, ok := d.GetOk("sensitivity"); ok {
		c.Sensitivity = attr.(int)
	}
	if attr, ok := d.GetOk("threshold"); ok {
		c.Threshold = attr.(int)
	}
	if attr, ok := d.GetOk("script"); ok {
		c.Script = attr.(string)
	}

}

func (TransactionCheck) setSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	d.Set("interval", c.Interval)
	d.Set("locations", c.Locations)
	d.Set("sensitivity", c.Sensitivity)
	d.Set("threshold", c.Threshold)
	d.Set("script", c.Script)
}

var transactionCheck TransactionCheck
