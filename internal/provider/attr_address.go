package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func AddressHostnameSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required: true,
		Description: `A valid DNS hostname (e.g., 'example.com', 'sub.example.com'). 
Must start and end with alphanumeric characters, can contain hyphens but not at the start or end, 
and must have at least one dot separator between valid DNS labels.`,
		Validators: []validator.String{
			HostnameValidator(),
		},
	}
}

func AddressHostnameOrIPSchemaAttribute() schema.StringAttribute {
	// Use this attribute when the address can be either a hostname or an IP (v4 or v6) address
	// Validation in that case done by API.
	return schema.StringAttribute{
		Description: `A valid DNS hostname (e.g., 'example.com') or IP address (IPv4 or IPv6). 
For hostnames: must start and end with alphanumeric characters. For IP addresses: supports 
both IPv4 (e.g., '192.168.1.1') and IPv6 (e.g., '2001:db8::1') formats.`,
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
		Description: `A valid URL with a required scheme (e.g., 'https://example.com', 'http://192.168.1.1:8080').
Must include protocol scheme and valid hostname or IP address. Port numbers are optional.`,
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
