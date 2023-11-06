package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func UseIPVersionSchemaAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Whether to use IPv4 or IPv6 for the check.",
		Default:     stringdefault.StaticString(""),
		Validators: []validator.String{
			OneOfStringValidator([]string{"", "IPV4", "IPV6"}),
		},
	}
}
