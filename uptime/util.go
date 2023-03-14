package uptime

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Takes the result of flatmap.Expand for a slice of strings
// and returns a list
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

// Expand a schema.Set into a slice of strings
func expandSetAttr(set interface{}) []string {
	return expandStringList(set.(*schema.Set).List())
}

func accumulateErrors(errs ...error) (err error) {
	for i := range errs {
		if errs[i] != nil {
			err = multierror.Append(err, errs[i])
		}
	}
	return err
}
