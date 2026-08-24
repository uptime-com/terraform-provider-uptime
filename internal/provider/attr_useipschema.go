package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func UseIPVersionSchemaAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		Description: "Internet Protocol version to use for the check. Valid values are \"IPV4\", \"IPV6\", " +
			"and \"\" for Any. Defaults to \"\".",
		Default: stringdefault.StaticString(""),
		Validators: []validator.String{
			OneOfStringValidator([]string{"", "IPV4", "IPV6"}),
		},
	}
}
