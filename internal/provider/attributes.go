package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//
// The most common attributes

func IDAttribute() schema.Int64Attribute {
	return schema.Int64Attribute{
		Computed: true,
	}
}

func URLAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
	}
}

func NameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required: true,
	}
}

//
// Less common attributes alphabetically

func ContactGroupsAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Default:     setdefault.StaticValue(mustDiag(types.SetValue(types.StringType, []attr.Value{types.StringValue("Default")}))),
	}
}

func IncludeInGlobalMetricsAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Include this check in uptime/response time calculations for the dashboard and status pages",
	}
}

func IntervalAttribute(defaultVal int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(defaultVal),
		Description: "The interval between checks in minutes",
	}
}

func IsPausedAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default:  booldefault.StaticBool(false),
	}
}

func LocationsAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Default: setdefault.StaticValue(mustDiag(types.SetValue(types.StringType,
			[]attr.Value{
				types.StringValue("US East"),
				types.StringValue("US West"),
			},
		))),
	}
}

func NotesAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
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

func SensitivityAttribute(defaultVal int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(defaultVal),
		Description: "How many locations should be down before an alert is sent",
	}
}

func TagsAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		Default:     setdefault.StaticValue(mustDiag(types.SetValue(types.StringType, []attr.Value{}))),
	}
}

func ThresholdAttribute(defaultVal int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: "A timeout alert will be issued if the check takes longer than this many seconds to complete",
		Default:     int64default.StaticInt64(defaultVal),
	}
}
