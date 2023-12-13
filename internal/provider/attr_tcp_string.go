package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func StringToSendSchemaAttribute() schema.Attribute {
	return schema.StringAttribute{
		Description: "String to send to the server",
		Optional:    true,
		Computed:    true,
	}
}

func StringToExpectSchemaAttribute() schema.Attribute {
	return schema.StringAttribute{
		Description: "String to expect in server response (may be repeated)",
		Optional:    true,
		Computed:    true,
	}
}
