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

func AddressHostnameOrIPSchemaAttribute() schema.StringAttribute {
	// Use this attribute when the address can be either a hostname or an IP (v4 or v6) address
	// Validation in that case done by API.
	return schema.StringAttribute{
		Required: true,
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
