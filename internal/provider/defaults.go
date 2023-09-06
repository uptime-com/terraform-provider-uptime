package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DefaultEmptyStringSet struct {
}

func (d *DefaultEmptyStringSet) Description(ctx context.Context) string {
	return "empty"
}

func (d *DefaultEmptyStringSet) MarkdownDescription(ctx context.Context) string {
	return "empty"
}

func (d *DefaultEmptyStringSet) DefaultSet(_ context.Context, _ defaults.SetRequest, rs *defaults.SetResponse) {
	rs.PlanValue = mustDiag(types.SetValue(types.StringType, []attr.Value{}))
}
