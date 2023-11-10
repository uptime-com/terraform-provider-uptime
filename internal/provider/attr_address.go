package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func AddressHostnameSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			HostnameValidator(),
		},
	}
}

func AddressHostnameSchemaAttributeDescription(desc string) schema.StringAttribute {
	return schema.StringAttribute{
		Description: desc,
		Required:    true,
		Validators: []validator.String{
			HostnameValidator(),
		},
	}
}

func AddressURLSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			URLValidator(),
		},
	}
}

func AddressURLSchemaAttributeDescription(desc string) schema.StringAttribute {
	return schema.StringAttribute{
		Description: desc,
		Required:    true,
		Validators: []validator.String{
			URLValidator(),
		},
	}
}
