package uptime

import (
	"context"
	"log"
	"net/http"
	"strconv"

	uptime "github.com/uptime-com/rest-api-clients/golang/uptime"
	"github.com/hashicorp/terraform/helper/schema"
)

// CheckType is an interface specifying the required methods for a valid parameter
// to the checkBuildFunc, checkCreateFunc, checkReadFunc, checkUpdateFunc, and
// checkDeleteFunc higher-order functions.
//
// The higher-order functions in this file allow for generalized creation of
// terraform resource schema CRUD functions.
type CheckType interface {
	typeStr() string
	getSpecificAttrs(*schema.ResourceData, *uptime.Check)
	setSpecificAttrs(*schema.ResourceData, *uptime.Check)
}

func setCommonCheckAttrs(d *schema.ResourceData, c *uptime.Check) {
	d.Set("address", c.Address)
	d.Set("contact_groups", c.ContactGroups)
	d.Set("name", c.Name)
	d.Set("notes", c.Notes)
	d.Set("tags", c.Tags)
}

func getCommonCheckAttrs(d *schema.ResourceData, c *uptime.Check) {
	if attr, ok := d.GetOk("name"); ok {
		c.Name = attr.(string)
	}

	if attr, ok := d.GetOk("tags"); ok {
		c.Tags = expandSetAttr(attr)
	}

	if attr, ok := d.GetOk("notes"); ok {
		c.Notes = attr.(string)
	}

	if attr, ok := d.GetOk("include_in_global_metrics"); ok {
		c.IncludeInGlobalMetrics = attr.(bool)
	}

}

func checkBuildFunc(ct CheckType) (func (d *schema.ResourceData) *uptime.Check) {
	return func(d *schema.ResourceData) *uptime.Check {
		check := &uptime.Check{
			CheckType: ct.typeStr(),
			Address: d.Get("address").(string),
			ContactGroups: expandSetAttr(d.Get("contact_groups")),
		}

		getCommonCheckAttrs(d, check)
		ct.getSpecificAttrs(d, check)

		return check
	}
}

func checkCreateFunc(ct CheckType) (schema.CreateFunc) {
	buildFunc := checkBuildFunc(ct)

	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*uptime.Client)
		ctx := context.Background()

		check := buildFunc(d)

		log.Printf("[INFO] Creating Uptime.com %s check for: %s", check.CheckType, check.Address)

		check, _, err := client.Checks.Create(ctx, check)
		if err != nil {
			return err
		}

		setResourceIDFromCheck(d, check)

		return checkReadFunc(ct)(d, meta)
	}
}

func checkReadFunc(ct CheckType) (schema.ReadFunc) {
	typeStr := ct.typeStr()

	return func (d *schema.ResourceData, meta interface{}) error {

		client := meta.(*uptime.Client)
		ctx := context.Background()

		log.Printf("[INFO] Reading Uptime.com %s check: %s", typeStr, d.Id())

		pk := pkFromResourceData(d)
		check, _, err := client.Checks.Get(ctx, pk)
		if err != nil {
			if uptErr, ok := err.(*uptime.Error); ok {
				if uptErr.Response.StatusCode == http.StatusNotFound {
					log.Printf("[WARN] Removing check %d from state because it no longer exists in Uptime.com", pk)
					d.SetId("")
					return nil
				}
			}
			return err
		}

		setCommonCheckAttrs(d, check)
		ct.setSpecificAttrs(d, check)

		return nil
	}

}

func checkUpdateFunc (ct CheckType) (schema.UpdateFunc) {
	buildFunc := checkBuildFunc(ct)

	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*uptime.Client)
		ctx := context.Background()

		check := buildFunc(d)

		pk := pkFromResourceData(d)
		check.PK = pk

		log.Printf("[DEBUG] Updating Domain %s check: %s", check.CheckType, d.Id())

		newCheck, _, err := client.Checks.Update(ctx, check)
		if err != nil {
			return err
		}

		setResourceIDFromCheck(d, newCheck)

		return checkReadFunc(ct)(d, meta)
	}
}

func checkDeleteFunc(ct CheckType) (schema.DeleteFunc) {
	typeStr := ct.typeStr()

	return func (d *schema.ResourceData, meta interface{}) error {
		client := meta.(*uptime.Client)
		ctx := context.Background()

		log.Printf("[INFO] Deleting Uptime %s check: %s", typeStr, d.Id())

		pk := pkFromResourceData(d)
		if _, err := client.Checks.Delete(ctx, pk); err != nil {
			return err
		}

		d.SetId("")
		return nil
	}
}

func setResourceIDFromCheck(d *schema.ResourceData, c *uptime.Check) {
	id := strconv.Itoa(c.PK)
	d.SetId(id)
}

func pkFromResourceData(d *schema.ResourceData) int {
	pk, _ := strconv.Atoi(d.Id())
	return pk
}
