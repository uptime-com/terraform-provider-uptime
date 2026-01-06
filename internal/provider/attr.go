package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

func IDSchemaAttribute() schema.Int64Attribute {
	return schema.Int64Attribute{
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	}
}

// ComputedIDSchemaAttribute returns an ID attribute without UseStateForUnknown.
// Use this for:
// - Nested object IDs (where new items get new IDs)
// - Resources that implement Update as Delete+Create (ID changes on update)
func ComputedIDSchemaAttribute() schema.Int64Attribute {
	return schema.Int64Attribute{
		Computed: true,
	}
}

func URLSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
	}
}

func NameSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required: true,
	}
}

func IncludeInGlobalMetricsSchemaAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Include this check in uptime/response time calculations for the dashboard and status pages",
	}
}

func IntervalSchemaAttribute(defaultVal int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(defaultVal),
		Description: "The interval between checks in minutes",
	}
}

func IsPausedSchemaAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default:  booldefault.StaticBool(false),
	}
}

func NotesSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("Managed by Terraform"),
	}
}

func NumRetriesAttribute(defaultVal int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(defaultVal),
		Description: "How many times the check should be retried before a location is considered down",
	}
}

func SensitivitySchemaAttribute(defaultVal int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(defaultVal),
		Description: "How many locations should be down before an alert is sent",
	}
}
