package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckGroupsDataSource(_ context.Context, p *providerImpl) datasource.DataSource {
	return CheckGroupsDataSource{p: p}
}

// CheckGroupsDataSchema defines the schema for the check_groups data source.
var CheckGroupsDataSchema = schema.Schema{
	Description: "Retrieve a list of all check groups configured in your Uptime.com account. Check groups combine multiple checks into a single logical unit.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Placeholder identifier for the data source",
		},
		"check_groups": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of all check groups in the account",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:    true,
						Description: "Unique identifier for the check group",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Human-readable name for the check group",
					},
					"contact_groups": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of contact group names assigned to this check group",
					},
					"tags": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "List of tags assigned to this check group",
					},
					"is_paused": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the check group is currently paused",
					},
					"notes": schema.StringAttribute{
						Computed:    true,
						Description: "Notes or description for the check group",
					},
					"include_in_global_metrics": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to include this check group in global metrics",
					},
					"sla": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "SLA configuration for the check group",
						Attributes: map[string]schema.Attribute{
							"uptime": schema.Float64Attribute{
								Computed:    true,
								Description: "Uptime SLA percentage target",
							},
							"latency": schema.Float64Attribute{
								Computed:    true,
								Description: "Response time SLA target in seconds",
							},
						},
					},
					"config": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Configuration for the check group",
						Attributes: map[string]schema.Attribute{
							"services": schema.ListAttribute{
								Computed:    true,
								ElementType: types.StringType,
								Description: "List of check names included in this group",
							},
							"tags": schema.ListAttribute{
								Computed:    true,
								ElementType: types.StringType,
								Description: "List of tags used to filter checks for this group",
							},
							"down_condition": schema.StringAttribute{
								Computed:    true,
								Description: "Condition that determines when the group check is considered DOWN",
							},
							"uptime_percent_calculation": schema.StringAttribute{
								Computed:    true,
								Description: "Method used to calculate the group's uptime percentage",
							},
							"response_time": schema.SingleNestedAttribute{
								Computed:    true,
								Description: "Response time calculation settings",
								Attributes: map[string]schema.Attribute{
									"calculation_mode": schema.StringAttribute{
										Computed:    true,
										Description: "Response time calculation mode",
									},
									"check_type": schema.StringAttribute{
										Computed:    true,
										Description: "Check type for response time calculation",
									},
									"single_check": schema.StringAttribute{
										Computed:    true,
										Description: "Single check for response time calculation",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

type CheckGroupsDataSourceModel struct {
	ID          types.String                     `tfsdk:"id"`
	CheckGroups []CheckGroupsDataSourceItemModel `tfsdk:"check_groups"`
}

type CheckGroupsDataSourceItemModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ContactGroups          types.List   `tfsdk:"contact_groups"`
	Tags                   types.List   `tfsdk:"tags"`
	IsPaused               types.Bool   `tfsdk:"is_paused"`
	Notes                  types.String `tfsdk:"notes"`
	IncludeInGlobalMetrics types.Bool   `tfsdk:"include_in_global_metrics"`
	SLA                    types.Object `tfsdk:"sla"`
	Config                 types.Object `tfsdk:"config"`
}

var _ datasource.DataSource = &CheckGroupsDataSource{}

type CheckGroupsDataSource struct {
	p *providerImpl
}

func (d CheckGroupsDataSource) Metadata(_ context.Context, rq datasource.MetadataRequest, rs *datasource.MetadataResponse) {
	rs.TypeName = rq.ProviderTypeName + "_check_groups"
}

func (d CheckGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, rs *datasource.SchemaResponse) {
	rs.Schema = CheckGroupsDataSchema
}

func (d CheckGroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, rs *datasource.ReadResponse) {
	api, err := d.p.api.Checks().List(ctx, upapi.CheckListOptions{})
	if err != nil {
		rs.Diagnostics.AddError("API call failed", err.Error())
		return
	}

	// Filter to only include check groups (checks with GroupConfig)
	var checkGroups []upapi.Check
	for _, check := range api {
		if check.GroupConfig != nil {
			checkGroups = append(checkGroups, check)
		}
	}

	model := CheckGroupsDataSourceModel{
		ID:          types.StringValue(""),
		CheckGroups: make([]CheckGroupsDataSourceItemModel, len(checkGroups)),
	}

	slaAttrTypes := map[string]attr.Type{
		"uptime":  types.Float64Type,
		"latency": types.Float64Type,
	}

	responseTimeAttrTypes := map[string]attr.Type{
		"calculation_mode": types.StringType,
		"check_type":       types.StringType,
		"single_check":     types.StringType,
	}

	configAttrTypes := map[string]attr.Type{
		"services":                   types.ListType{ElemType: types.StringType},
		"tags":                       types.ListType{ElemType: types.StringType},
		"down_condition":             types.StringType,
		"uptime_percent_calculation": types.StringType,
		"response_time":              types.ObjectType{AttrTypes: responseTimeAttrTypes},
	}

	for i, check := range checkGroups {
		// Convert contact groups slice to types.List
		contactGroupsList := types.ListNull(types.StringType)
		if check.ContactGroups != nil && len(*check.ContactGroups) > 0 {
			contactGroupsValues := make([]attr.Value, len(*check.ContactGroups))
			for j, v := range *check.ContactGroups {
				contactGroupsValues[j] = types.StringValue(v)
			}
			contactGroupsList = types.ListValueMust(types.StringType, contactGroupsValues)
		}

		// Convert tags slice to types.List
		tagsList := types.ListNull(types.StringType)
		if len(check.Tags) > 0 {
			tagsValues := make([]attr.Value, len(check.Tags))
			for j, v := range check.Tags {
				tagsValues[j] = types.StringValue(v)
			}
			tagsList = types.ListValueMust(types.StringType, tagsValues)
		}

		// Build SLA object
		uptimeFloat, _ := check.UptimeSLA.Float64()
		latencyFloat, _ := check.ResponseTimeSLA.Float64()
		slaObj := types.ObjectValueMust(slaAttrTypes, map[string]attr.Value{
			"uptime":  types.Float64Value(uptimeFloat),
			"latency": types.Float64Value(latencyFloat),
		})

		// Build config object
		var configObj types.Object
		if check.GroupConfig != nil {
			// Convert services slice to types.List
			servicesList := types.ListNull(types.StringType)
			if len(check.GroupConfig.CheckServices) > 0 {
				servicesValues := make([]attr.Value, len(check.GroupConfig.CheckServices))
				for j, v := range check.GroupConfig.CheckServices {
					servicesValues[j] = types.StringValue(v)
				}
				servicesList = types.ListValueMust(types.StringType, servicesValues)
			}

			// Convert config tags slice to types.List
			configTagsList := types.ListNull(types.StringType)
			if len(check.GroupConfig.CheckTags) > 0 {
				configTagsValues := make([]attr.Value, len(check.GroupConfig.CheckTags))
				for j, v := range check.GroupConfig.CheckTags {
					configTagsValues[j] = types.StringValue(v)
				}
				configTagsList = types.ListValueMust(types.StringType, configTagsValues)
			}

			// Build response_time object
			responseTimeObj := types.ObjectValueMust(responseTimeAttrTypes, map[string]attr.Value{
				"calculation_mode": types.StringValue(check.GroupConfig.ResponseTimeCalculationMode),
				"check_type":       types.StringValue(check.GroupConfig.ResponseTimeCheckType),
				"single_check":     types.StringValue(check.GroupConfig.ResponseTimeSingleCheck),
			})

			configObj = types.ObjectValueMust(configAttrTypes, map[string]attr.Value{
				"services":                   servicesList,
				"tags":                       configTagsList,
				"down_condition":             types.StringValue(check.GroupConfig.CheckDownCondition),
				"uptime_percent_calculation": types.StringValue(check.GroupConfig.UptimePercentCalculation),
				"response_time":              responseTimeObj,
			})
		} else {
			configObj = types.ObjectNull(configAttrTypes)
		}

		model.CheckGroups[i] = CheckGroupsDataSourceItemModel{
			ID:                     types.Int64Value(check.PK),
			Name:                   types.StringValue(check.Name),
			ContactGroups:          contactGroupsList,
			Tags:                   tagsList,
			IsPaused:               types.BoolValue(check.IsPaused),
			Notes:                  types.StringValue(check.Notes),
			IncludeInGlobalMetrics: types.BoolValue(check.IncludeInGlobalMetrics),
			SLA:                    slaObj,
			Config:                 configObj,
		}
	}

	rs.Diagnostics.Append(rs.State.Set(ctx, model)...)
}
