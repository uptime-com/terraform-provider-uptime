package uptime

import (
	"github.com/hashicorp/terraform/helper/schema"
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
