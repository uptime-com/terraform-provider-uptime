package uptime

import (
	"fmt"

	uptime "github.com/uptime-com/rest-api-clients/golang/uptime"
	"github.com/hashicorp/terraform/helper/schema"

	"strings"
)

func resourceUptimeCheckHTTP() *schema.Resource {
	return &schema.Resource{
		Create: checkCreateFunc(httpCheck),
		Read: checkReadFunc(httpCheck),
		Update: checkUpdateFunc(httpCheck),
		Delete: checkDeleteFunc(httpCheck),
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
			"expect_string": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"send_string": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
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
			"username": {
				Type: schema.TypeString,
				Optional: true,
			},
			"password": {
				Type: schema.TypeString,
				Optional: true,
			},
			"headers": {
				Type: schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"port": {
				Type: schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// HTTPCheck implements the CheckType interface for Uptime.com HTTP checks.
type HTTPCheck struct{}

func (HTTPCheck) typeStr() string {return "HTTP"}

func (HTTPCheck) getSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	if attr, ok := d.GetOk("interval"); ok {
		c.Interval = attr.(int)
	}
	if attr, ok := d.GetOk("locations"); ok {
		c.Locations = expandSetAttr(attr)
	}
	if attr, ok := d.GetOk("ip_version"); ok {
		c.IPVersion = attr.(string)
	}
	if attr, ok := d.GetOk("expect_string"); ok {
		c.ExpectString = attr.(string)
	}
	if attr, ok := d.GetOk("send_string"); ok {
		c.SendString = attr.(string)
	}
	if attr, ok := d.GetOk("sensitivity"); ok {
		c.Sensitivity = attr.(int)
	}
	if attr, ok := d.GetOk("threshold"); ok {
		c.Threshold = attr.(int)
	}
	if attr, ok := d.GetOk("username"); ok {
		c.Username = attr.(string)
	}
	if attr, ok := d.GetOk("password"); ok {
		c.Password = attr.(string)
	}
	if attr, ok := d.GetOk("port"); ok {
		c.Port = attr.(int)
	}
	if attr, ok := d.GetOk("headers"); ok {
		hs := headersMapToString(attr.(map[string]interface{}))
		c.Headers = hs
	}
}

func (HTTPCheck) setSpecificAttrs(d *schema.ResourceData, c *uptime.Check) {
	d.Set("interval", c.Interval)
	d.Set("locations", c.Locations)
	d.Set("ip_version", c.IPVersion)
	d.Set("expect_string", c.ExpectString)
	d.Set("send_string", c.SendString)
	d.Set("sensitivity", c.Sensitivity)
	d.Set("threshold", c.Threshold)
	d.Set("username", c.Username)
	d.Set("password", c.Password)
	d.Set("port", c.Port)

	hs := headersStringToMap(c.Headers)
	d.Set("headers", hs)
}

var httpCheck HTTPCheck

func headersStringToMap(hs string) map[string]string {
	hm := make(map[string]string)
	if hs != "" {
		headers := strings.Split(hs, "\n")
		for _, h := range headers {
			kv := strings.Split(h, ": ")
			if len(kv) > 1 {
				hm[kv[0]] = kv[1]
			}
		}
	}
	return hm
}

func headersMapToString(hm map[string]interface{}) string {
	var hs strings.Builder

	for k, v := range hm {
		s := fmt.Sprintf("%s: %s\n", k, v.(string))
		hs.WriteString(s)
	}
	return hs.String()

}
