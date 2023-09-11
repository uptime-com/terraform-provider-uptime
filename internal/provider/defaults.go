package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringSetDefault(v []string) defaults.Set {
	return &stringSetDefault{
		v: v,
	}
}

func StringSetEmptyDefault() defaults.Set {
	return StringSetDefault(nil)
}

type stringSetDefault struct {
	v []string
}

func (d *stringSetDefault) Description(ctx context.Context) string {
	return ""
}

func (d *stringSetDefault) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (d *stringSetDefault) DefaultSet(_ context.Context, _ defaults.SetRequest, rs *defaults.SetResponse) {
	var v []attr.Value
	for i := range d.v {
		v = append(v, types.StringValue(d.v[i]))
	}
	rs.PlanValue = types.SetValueMust(types.StringType, v)
}
