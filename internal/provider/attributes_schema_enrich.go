package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func EnreachResourceSchema(s schema.Schema, p *providerImpl) schema.Schema {
	if locations, ok := s.Attributes["locations"]; ok {
		if attr, ok := locations.(schema.SetAttribute); ok {
			attr.Validators = append(attr.Validators,
				setvalidator.ValueStringsAre(stringvalidator.OneOf(p.locations...)))
		}
	}
	return s
}
