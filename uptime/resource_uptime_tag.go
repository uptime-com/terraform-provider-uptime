package uptime

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/uptime-com/uptime-client-go"
)

func resourceUptimeCheckTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceUptimeCheckTagCreate,
		Read:   resourceUptimeCheckTagRead,
		Update: resourceUptimeCheckTagUpdate,
		Delete: resourceUptimeCheckTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUptimeCheckTagImport,
		},
		Schema: map[string]*schema.Schema{
			"tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"color_hex": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildUptimeCheckTag(d *schema.ResourceData) *uptime.Tag {
	checkTag := &uptime.Tag{
		Tag:      d.Get("tag").(string),
		ColorHex: d.Get("color_hex").(string),
		URL:      d.Get("url").(string),
	}
	return checkTag
}

func resourceUptimeCheckTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*uptime.Client)
	ctx := context.Background()

	t := buildUptimeCheckTag(d)

	log.Printf("[INFO] Creating Uptime.com check tag: %s with color %s", t.Tag, t.ColorHex)
	t, _, err := client.Tags.Create(ctx, t)
	if err != nil {
		return err
	}
	setResourceIDFromTag(d, t)

	return resourceUptimeCheckTagRead(d, meta)
}

func resourceUptimeCheckTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*uptime.Client)
	ctx := context.Background()

	log.Printf("[INFO] Reading Uptime.com check tag: %s", d.Id())

	pk := pkFromResourceData(d)
	t, _, err := client.Tags.Get(ctx, pk)
	if err != nil {
		if uptErr, ok := err.(*uptime.Error); ok {
			if uptErr.Response.StatusCode == http.StatusNotFound {
				log.Printf("[WARN] Removing tag %d from state because it no longer exists in Uptime.com", pk)
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("tag", t.Tag)
	d.Set("color_hex", t.ColorHex)
	d.Set("url", t.URL)
	return nil
}

func resourceUptimeCheckTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*uptime.Client)
	ctx := context.Background()

	t := buildUptimeCheckTag(d)
	pk := pkFromResourceData(d)
	t.PK = pk

	log.Printf("[DEBUG] Updating tag: %s", d.Id())

	newCheck, _, err := client.Tags.Update(ctx, t)
	if err != nil {
		return err
	}

	setResourceIDFromTag(d, newCheck)

	return resourceUptimeCheckTagRead(d, meta)
}

func resourceUptimeCheckTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*uptime.Client)
	ctx := context.Background()

	log.Printf("[INFO] Deleting Uptime tag: %s", d.Id())

	pk := pkFromResourceData(d)
	if _, err := client.Tags.Delete(ctx, pk); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceUptimeCheckTagImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*uptime.Client)

	log.Printf("[INFO] Importing Uptime tag: %s", d.Id())

	pk := pkFromResourceData(d)
	tag, res, err := client.Tags.Get(ctx, pk)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	e := multierror.Append(new(multierror.Error),
		d.Set("tag", tag.Tag),
		d.Set("color_hex", tag.ColorHex),
		d.Set("url", tag.URL),
	)
	return []*schema.ResourceData{d}, e.ErrorOrNil()
}

func setResourceIDFromTag(d *schema.ResourceData, t *uptime.Tag) {
	id := strconv.Itoa(t.PK)
	d.SetId(id)
}
