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

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"
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

// SLA Attributes

func ResponseTimeSLAAttribute(defaultVal string) schema.StringAttribute {
	return schema.StringAttribute{
		CustomType:  customtypes.DurationType{},
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(defaultVal),
		Description: "The maximum average response time. Unit is mandatory (e.g. 1500ms or 1.5s or 1s500ms).",
	}
}

//
// Less common attributes alphabetically

func AddressHostnameAttribute() schema.StringAttribute {
	return AddressHostnameAttributeDesc("The hostname to check")
}

func AddressHostnameAttributeDesc(desc string) schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: desc,
		Validators: []validator.String{
			HostnameValidator(),
		},
	}
}

func AddressURLAttribute() schema.StringAttribute {
	return AddressURLAttributeDesc("The URL to check")
}

func AddressURLAttributeDesc(desc string) schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: desc,
		Validators: []validator.String{
			URLValidator(),
		},
	}
}

func ContactGroupsAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Default")})),
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
		Default: setdefault.StaticValue(types.SetValueMust(types.StringType,
			[]attr.Value{
				types.StringValue("US East"),
				types.StringValue("US West"),
			},
		)),
	}
}

func LocationsReadOnlyAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
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
		Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
	}
}

func ThresholdAttribute(defaultVal int64) schema.Int64Attribute {
	return ThresholdDescriptionAttribute(defaultVal, "A timeout alert will be issued if the check takes longer than this many seconds to complete")
}

func ThresholdDescriptionAttribute(defaultVal int64, desc string) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: desc,
		Default:     int64default.StaticInt64(defaultVal),
	}
}

func UseIPVersionAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			OneOfStringValidator([]string{"", "IPV4", "IPV6"}),
		},
	}
}
