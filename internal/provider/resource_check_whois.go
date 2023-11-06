package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckWHOISResource(_ context.Context, p *providerImpl) resource.Resource {
	return APIResource[CheckWHOISResourceModel, upapi.CheckWHOIS, upapi.Check]{
		api: CheckWHOISResourceAPI{provider: p},
		mod: CheckWHOISResourceModelAdapter{},
		meta: APIResourceMetadata{
			TypeNameSuffix: "check_whois",
			Schema:         CheckWHOISResourceSchema,
		},
	}
}

var CheckWHOISResourceSchema = schema.Schema{
	Description: "Monitor domain's expiry date and registration details",
	Attributes: map[string]schema.Attribute{
		"id":             IDSchemaAttribute(),
		"url":            URLSchemaAttribute(),
		"name":           NameSchemaAttribute(),
		"contact_groups": ContactGroupsSchemaAttribute(),
		"locations":      LocationsReadOnlySchemaAttribute(),
		"tags":           TagsSchemaAttribute(),
		"is_paused":      IsPausedSchemaAttribute(),
		"threshold":      ThresholdDescriptionSchemaAttribute(20, "Raise an alert if there are less than this many days before the domain needs to be renewed."),
		"num_retries":    NumRetriesSchemaAttribute(2),
		"notes":          NotesSchemaAttribute(),
		"address":        AddressHostnameSchemaAttribute(),
		"expect_string": schema.StringAttribute{
			Required:    true,
			Description: "The current domain registration info that should always match.",
		},

		// NOTE: for this check only uptime SLA is meaningful. Latency is ignored even if provided.
		// TODO: implement partial SLA attribute containing only uptime.
		"sla": SLASchemaAttribute(),
	},
}

type CheckWHOISResourceModel struct {
	ID            types.Int64   `tfsdk:"id"  ref:"PK,opt"`
	URL           types.String  `tfsdk:"url" ref:"URL,opt"`
	Name          types.String  `tfsdk:"name"`
	ContactGroups types.Set     `tfsdk:"contact_groups"`
	Locations     types.Set     `tfsdk:"locations"`
	Tags          types.Set     `tfsdk:"tags"`
	IsPaused      types.Bool    `tfsdk:"is_paused"`
	Address       types.String  `tfsdk:"address"`
	ExpectString  types.String  `tfsdk:"expect_string"`
	Threshold     types.Int64   `tfsdk:"threshold"`
	NumRetries    types.Int64   `tfsdk:"num_retries"`
	Notes         types.String  `tfsdk:"notes"`
	SLA           types.Object  `tfsdk:"sla"`
	sla           *SLAAttribute `tfsdk:"-"`
}

func (m CheckWHOISResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckWHOISResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckWHOISResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckWHOISResourceModel, diag.Diagnostics) {
	model := *new(CheckWHOISResourceModel)
	diags := sg.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}
	model.sla, diags = a.SLAAttributeContext(ctx, model.SLA)
	if diags.HasError() {
		return nil, diags
	}
	return &model, nil
}

func (a CheckWHOISResourceModelAdapter) ToAPIArgument(model CheckWHOISResourceModel) (*upapi.CheckWHOIS, error) {
	api := upapi.CheckWHOIS{
		Name:          model.Name.ValueString(),
		ContactGroups: a.ContactGroups(model.ContactGroups),
		Locations:     a.Locations(model.Locations),
		Tags:          a.Tags(model.Tags),
		IsPaused:      model.IsPaused.ValueBool(),
		Address:       model.Address.ValueString(),
		ExpectString:  model.ExpectString.ValueString(),
		Threshold:     model.Threshold.ValueInt64(),
		NumRetries:    model.NumRetries.ValueInt64(),
		Notes:         model.Notes.ValueString(),
	}
	if model.sla != nil {
		if !model.sla.Uptime.IsUnknown() {
			api.UptimeSLA = model.sla.Uptime.ValueDecimal()
		}
	}

	return &api, nil
}

func (a CheckWHOISResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckWHOISResourceModel, error) {
	model := CheckWHOISResourceModel{
		ID:            types.Int64Value(api.PK),
		URL:           types.StringValue(api.URL),
		Name:          types.StringValue(api.Name),
		ContactGroups: a.ContactGroupsValue(api.ContactGroups),
		Locations:     a.LocationsValue(api.Locations),
		Tags:          a.TagsValue(api.Tags),
		IsPaused:      types.BoolValue(api.IsPaused),
		Address:       types.StringValue(api.Address),
		ExpectString:  types.StringValue(api.ExpectString),
		Threshold:     types.Int64Value(api.Threshold),
		NumRetries:    types.Int64Value(api.NumRetries),
		Notes:         types.StringValue(api.Notes),
		SLA: a.SLAAttributeValue(SLAAttribute{
			Latency: DurationValue(0), // NOTE: latency is not meaningful for this check, using bogus value
			Uptime:  DecimalValue(api.UptimeSLA),
		}),
	}
	return &model, nil
}

type CheckWHOISResourceAPI struct {
	provider *providerImpl
}

func (c CheckWHOISResourceAPI) Create(ctx context.Context, arg upapi.CheckWHOIS) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateWHOIS(ctx, arg)
}

func (c CheckWHOISResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckWHOISResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckWHOIS) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateWHOIS(ctx, pk, arg)
}

func (c CheckWHOISResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
