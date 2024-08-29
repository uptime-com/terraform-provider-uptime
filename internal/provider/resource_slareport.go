package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewSLAReportResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[SLAReportResourceModel, upapi.SLAReport, upapi.SLAReport]{
		api: &SLAReportResourceAPI{provider: p},
		mod: SLAReportResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "sla_report",
			Schema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":            IDSchemaAttribute(),
					"url":           URLSchemaAttribute(),
					"name":          NameSchemaAttribute(),
					"services_tags": TagsSchemaAttribute(),
					"services_selected": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Optional: true,
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
							},
						},
					},
					"reporting_groups": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Computed: true,
									Optional: true,
								},
								"name": NameSchemaAttribute(),
								"group_services": schema.SetAttribute{
									Optional:    true,
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"default_date_range": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("TODAY"),
						Validators: []validator.String{
							OneOfStringValidator([]string{
								"TODAY", "YESTERDAY", "LAST_7D", "LAST_30D", "THIS_WEEK", "LAST_WEEK", "THIS_MONTH", "LAST_MONTH",
							}),
						},
					},
					"show_uptime_section": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"show_uptime_sla": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"show_response_time_section": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"show_response_time_sla": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"filter_with_downtime": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"filter_uptime_sla_violations": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"filter_response_time_sla_violations": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"filter_slowest": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"uptime_section_sort": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("BY_UPTIME"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"BY_UPTIME", "BY_SLA"}),
						},
					},
					"response_time_section_sort": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("BY_RESPONSE_TIME"),
						Validators: []validator.String{
							OneOfStringValidator([]string{"BY_RESPONSE_TIME", "BY_SLA"}),
						},
					},
				},
			},
		},
	}
}

type SLAReportResourceModel struct {
	ID                              types.Int64  `tfsdk:"id"`
	URL                             types.String `tfsdk:"url"`
	Name                            types.String `tfsdk:"name"`
	ServicesTags                    types.Set    `tfsdk:"services_tags"`
	SelectedServices                types.Set    `tfsdk:"services_selected"`
	ReportingGroups                 types.Set    `tfsdk:"reporting_groups"`
	DefaultDateRange                types.String `tfsdk:"default_date_range"`
	ShowUptimeSection               types.Bool   `tfsdk:"show_uptime_section"`
	ShowUptimeSLA                   types.Bool   `tfsdk:"show_uptime_sla"`
	ShowResponseTimeSection         types.Bool   `tfsdk:"show_response_time_section"`
	ShowResponseTimeSLA             types.Bool   `tfsdk:"show_response_time_sla"`
	FilterWithDowntime              types.Bool   `tfsdk:"filter_with_downtime"`
	FilterUptimeSLAViolations       types.Bool   `tfsdk:"filter_uptime_sla_violations"`
	FilterResponseTimeSLAViolations types.Bool   `tfsdk:"filter_response_time_sla_violations"`
	FilterSlowest                   types.Bool   `tfsdk:"filter_slowest"`
	UptimeSectionSort               types.String `tfsdk:"uptime_section_sort"`
	ResponseTimeSectionSort         types.String `tfsdk:"response_time_section_sort"`

	servicesSelected []ServicesSelectedAttribute `tfsdk:"-"`
	reportingGroups  []ReportingGroupsAttribute  `tfsdk:"-"`
}

func (m SLAReportResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type SLAReportResourceModelAdapter struct {
	TagsAttributeAdapter
	SetAttributeAdapter[string]
}

func (a SLAReportResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*SLAReportResourceModel, diag.Diagnostics) {
	model := *new(SLAReportResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}

	model.servicesSelected, diags = a.ServicesSelectedContext(ctx, model.SelectedServices)
	if diags.HasError() {
		return nil, diags
	}

	model.reportingGroups, diags = a.ReportingGroupsContext(ctx, model.ReportingGroups)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a SLAReportResourceModelAdapter) ToAPIArgument(model SLAReportResourceModel) (*upapi.SLAReport, error) {
	api := upapi.SLAReport{
		Name:                            model.Name.ValueString(),
		DefaultDateRange:                model.DefaultDateRange.ValueString(),
		ServicesTags:                    a.Tags(model.ServicesTags),
		ShowUptimeSection:               model.ShowUptimeSection.ValueBool(),
		ShowUptimeSLA:                   model.ShowUptimeSLA.ValueBool(),
		ShowResponseTimeSection:         model.ShowResponseTimeSection.ValueBool(),
		ShowResponseTimeSLA:             model.ShowResponseTimeSLA.ValueBool(),
		FilterWithDowntime:              model.FilterWithDowntime.ValueBool(),
		FilterUptimeSLAViolations:       model.FilterUptimeSLAViolations.ValueBool(),
		FilterResponseTimeSLAViolations: model.FilterResponseTimeSLAViolations.ValueBool(),
		FilterSlowest:                   model.FilterSlowest.ValueBool(),
		UptimeSectionSort:               model.UptimeSectionSort.ValueString(),
		ResponseTimeSectionSort:         model.ResponseTimeSectionSort.ValueString(),
	}

	if len(model.servicesSelected) != 0 {
		ss := make([]upapi.SLAReportService, 0)
		for _, v := range model.servicesSelected {
			ss = append(ss, upapi.SLAReportService{
				PK:   int(v.ID.ValueInt64()),
				Name: v.Name.ValueString(),
			})
		}
		api.ServicesSelected = &ss
	}

	if len(model.reportingGroups) != 0 {
		rg := make([]upapi.SLAReportGroup, 0)
		for _, v := range model.reportingGroups {
			rg = append(rg, upapi.SLAReportGroup{
				ID:            int(v.ID.ValueInt64()),
				Name:          v.Name.ValueString(),
				GroupServices: a.Slice(v.GroupServices),
			})
		}
		api.ReportingGroups = &rg
	}
	return &api, nil
}

func (a SLAReportResourceModelAdapter) FromAPIResult(api upapi.SLAReport) (*SLAReportResourceModel, error) {
	model := SLAReportResourceModel{
		ID:                              types.Int64Value(api.PK),
		URL:                             types.StringValue(api.URL),
		Name:                            types.StringValue(api.Name),
		ServicesTags:                    a.TagsValue(api.ServicesTags),
		DefaultDateRange:                types.StringValue(api.DefaultDateRange),
		ShowUptimeSection:               types.BoolValue(api.ShowUptimeSection),
		ShowUptimeSLA:                   types.BoolValue(api.ShowUptimeSLA),
		ShowResponseTimeSection:         types.BoolValue(api.ShowResponseTimeSection),
		ShowResponseTimeSLA:             types.BoolValue(api.ShowResponseTimeSLA),
		FilterWithDowntime:              types.BoolValue(api.FilterWithDowntime),
		FilterUptimeSLAViolations:       types.BoolValue(api.FilterUptimeSLAViolations),
		FilterResponseTimeSLAViolations: types.BoolValue(api.FilterResponseTimeSLAViolations),
		FilterSlowest:                   types.BoolValue(api.FilterSlowest),
		UptimeSectionSort:               types.StringValue(api.UptimeSectionSort),
		ResponseTimeSectionSort:         types.StringValue(api.ResponseTimeSectionSort),
	}

	var diags diag.Diagnostics
	servicesSelected := []ServicesSelectedAttribute{}
	for _, item := range *api.ServicesSelected {
		servicesSelected = append(servicesSelected, ServicesSelectedAttribute{
			ID:   types.Int64Value(int64(item.PK)),
			Name: types.StringValue(item.Name),
		})
	}

	if model.SelectedServices, diags = a.ServicesSelectedValue(servicesSelected); diags.HasError() {
		return nil, fmt.Errorf("failed to convert selected services: %v", diags)
	}

	reportingGroups := []ReportingGroupsAttribute{}
	for _, item := range *api.ReportingGroups {
		reportingGroups = append(reportingGroups, ReportingGroupsAttribute{
			ID:            types.Int64Value(int64(item.ID)),
			Name:          types.StringValue(item.Name),
			GroupServices: a.SliceValue(item.GroupServices),
		})
	}

	if model.ReportingGroups, diags = a.ReportingGroupsValue(reportingGroups); diags.HasError() {
		return nil, fmt.Errorf("failed to convert reporting groups: %v", diags)
	}

	return &model, nil
}

func (a SLAReportResourceModelAdapter) ServicesSelectedContext(ctx context.Context, v types.Set) ([]ServicesSelectedAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	out := make([]ServicesSelectedAttribute, 0)
	if d := v.ElementsAs(ctx, &out, false); d.HasError() {
		return nil, d
	}
	return out, nil
}

func (a SLAReportResourceModelAdapter) ServicesSelectedValue(model []ServicesSelectedAttribute) (types.Set, diag.Diagnostics) {
	values, diags := a.servicesAttributeValues(model)
	if diags.HasError() {
		return types.Set{}, diags
	}
	return types.SetValueMust(
		types.ObjectType{}.WithAttributeTypes(a.servicesAttributeTypes()), values), diags
}

func (a SLAReportResourceModelAdapter) servicesAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.Int64Type,
		"name": types.StringType,
	}
}

func (a SLAReportResourceModelAdapter) servicesAttributeValues(model []ServicesSelectedAttribute) (out []attr.Value, diags diag.Diagnostics) {
	out = make([]attr.Value, len(model))
	for i := range model {
		out[i], diags = types.ObjectValue(a.servicesAttributeTypes(), map[string]attr.Value{
			"id":   model[i].ID,
			"name": model[i].Name,
		})
		if diags.HasError() {
			return
		}
	}
	return
}

func (a SLAReportResourceModelAdapter) ReportingGroupsContext(ctx context.Context, v types.Set) ([]ReportingGroupsAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	out := make([]ReportingGroupsAttribute, 0)
	if d := v.ElementsAs(ctx, &out, false); d.HasError() {
		return nil, d
	}
	return out, nil
}

func (a SLAReportResourceModelAdapter) ReportingGroupsValue(model []ReportingGroupsAttribute) (types.Set, diag.Diagnostics) {
	values, diags := a.reportingGroupsAttributeValues(model)
	if diags.HasError() {
		return types.Set{}, diags
	}
	return types.SetValueMust(
		types.ObjectType{}.WithAttributeTypes(a.reportingGroupsAttributeTypes()), values), diags
}

func (a SLAReportResourceModelAdapter) reportingGroupsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.Int64Type,
		"name":           types.StringType,
		"group_services": types.SetType{}.WithElementType(types.StringType),
	}
}

func (a SLAReportResourceModelAdapter) reportingGroupsAttributeValues(model []ReportingGroupsAttribute) (out []attr.Value, diags diag.Diagnostics) {
	out = make([]attr.Value, len(model))
	for i := range model {
		out[i], diags = types.ObjectValue(a.reportingGroupsAttributeTypes(), map[string]attr.Value{
			"id":             model[i].ID,
			"name":           model[i].Name,
			"group_services": model[i].GroupServices,
		})
		if diags.HasError() {
			return
		}
	}
	return
}

type SLAReportResourceAPI struct {
	provider *providerImpl
}

func (a SLAReportResourceAPI) Create(ctx context.Context, arg upapi.SLAReport) (*upapi.SLAReport, error) {
	rg := arg.ReportingGroups

	arg.ReportingGroups = nil
	obj, err := a.provider.api.SLAReports().Create(ctx, arg)
	if err != nil {
		return nil, err
	}

	if rg != nil {
		if err := a.createGroups(ctx, obj, *rg); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func (a SLAReportResourceAPI) Read(ctx context.Context, arg upapi.PrimaryKeyable) (*upapi.SLAReport, error) {
	obj, err := a.provider.api.SLAReports().Get(ctx, arg)
	return obj, err
}

func (a SLAReportResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.SLAReport) (*upapi.SLAReport, error) {
	rg := arg.ReportingGroups

	arg.ReportingGroups = nil
	obj, err := a.provider.api.SLAReports().Update(ctx, pk, arg)
	if err != nil {
		return nil, err
	}

	if rg != nil {
		if err := a.createGroups(ctx, obj, *rg); err != nil {
			return nil, err
		}
	}
	return obj, err
}

func (a SLAReportResourceAPI) Delete(ctx context.Context, keyable upapi.PrimaryKeyable) error {
	return a.provider.api.SLAReports().Delete(ctx, keyable)
}

func (a SLAReportResourceAPI) createGroups(ctx context.Context, obj *upapi.SLAReport, rgs []upapi.SLAReportGroup) error {
	rgObjList := make([]upapi.SLAReportGroup, 0)
	for _, rg := range rgs {
		rgObj, err := a.provider.api.SLAReports().ReportingGroups(obj).Create(ctx, rg)
		if err != nil {
			return err
		}
		rgObjList = append(rgObjList, *rgObj)
	}
	obj.ReportingGroups = &rgObjList
	return nil
}

type ServicesSelectedAttribute struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type ReportingGroupsAttribute struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	GroupServices types.Set    `tfsdk:"group_services"`
}
