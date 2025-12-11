package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewCheckRDAPResource(_ context.Context, p *providerImpl) resource.Resource {
	return NewImportableAPIResource[CheckRDAPResourceModel, upapi.CheckRDAP, upapi.Check](
		CheckRDAPResourceAPI{provider: p},
		CheckRDAPResourceModelAdapter{},
		APIResourceMetadata{
			TypeNameSuffix: "check_rdap",
			Schema: schema.Schema{
				Description: "Monitor domain's expiry date and registration details using RDAP (Registration Data Access Protocol). Import using the check ID: `terraform import uptime_check_rdap.example 123`",
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
					"expect_string":  StringToExpectSchemaAttribute(),
					"send_resolved_notifications": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to send notifications when the check recovers from a down state.",
					},
					"sla": SLASchemaAttribute(),
				},
			},
		},
		ImportStateSimpleID,
	)
}

type CheckRDAPResourceModel struct {
	ID                        types.Int64   `tfsdk:"id"  ref:"PK,opt"`
	URL                       types.String  `tfsdk:"url" ref:"URL,opt"`
	Name                      types.String  `tfsdk:"name"`
	ContactGroups             types.Set     `tfsdk:"contact_groups"`
	Locations                 types.Set     `tfsdk:"locations"`
	Tags                      types.Set     `tfsdk:"tags"`
	IsPaused                  types.Bool    `tfsdk:"is_paused"`
	Address                   types.String  `tfsdk:"address"`
	ExpectString              types.String  `tfsdk:"expect_string"`
	Threshold                 types.Int64   `tfsdk:"threshold"`
	NumRetries                types.Int64   `tfsdk:"num_retries"`
	Notes                     types.String  `tfsdk:"notes"`
	SendResolvedNotifications types.Bool    `tfsdk:"send_resolved_notifications"`
	SLA                       types.Object  `tfsdk:"sla"`
	sla                       *SLAAttribute `tfsdk:"-"`
}

func (m CheckRDAPResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type CheckRDAPResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
	LocationsAttributeAdapter
	TagsAttributeAdapter
	SLAAttributeContextAdapter
}

func (a CheckRDAPResourceModelAdapter) Get(ctx context.Context, sg StateGetter) (*CheckRDAPResourceModel, diag.Diagnostics) {
	model := *new(CheckRDAPResourceModel)
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

func (a CheckRDAPResourceModelAdapter) ToAPIArgument(model CheckRDAPResourceModel) (*upapi.CheckRDAP, error) {
	api := upapi.CheckRDAP{
		Name:                      model.Name.ValueString(),
		ContactGroups:             a.ContactGroups(model.ContactGroups),
		Locations:                 a.Locations(model.Locations),
		Tags:                      a.Tags(model.Tags),
		IsPaused:                  model.IsPaused.ValueBool(),
		Address:                   model.Address.ValueString(),
		ExpectString:              model.ExpectString.ValueString(),
		Threshold:                 model.Threshold.ValueInt64(),
		NumRetries:                model.NumRetries.ValueInt64(),
		Notes:                     model.Notes.ValueString(),
		SendResolvedNotifications: model.SendResolvedNotifications.ValueBool(),
	}
	if model.sla != nil {
		if !model.sla.Uptime.IsUnknown() {
			api.UptimeSLA = model.sla.Uptime.ValueDecimal()
		}
	}

	return &api, nil
}

func (a CheckRDAPResourceModelAdapter) FromAPIResult(api upapi.Check) (*CheckRDAPResourceModel, error) {
	model := CheckRDAPResourceModel{
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
			Latency: DurationValue(0),
			Uptime:  DecimalValue(api.UptimeSLA),
		}),
	}
	return &model, nil
}

type CheckRDAPResourceAPI struct {
	provider *providerImpl
}

func (c CheckRDAPResourceAPI) Create(ctx context.Context, arg upapi.CheckRDAP) (*upapi.Check, error) {
	return c.provider.api.Checks().CreateRDAP(ctx, arg)
}

func (c CheckRDAPResourceAPI) Read(ctx context.Context, pk upapi.PrimaryKeyable) (*upapi.Check, error) {
	return c.provider.api.Checks().Get(ctx, pk)
}

func (c CheckRDAPResourceAPI) Update(ctx context.Context, pk upapi.PrimaryKeyable, arg upapi.CheckRDAP) (*upapi.Check, error) {
	return c.provider.api.Checks().UpdateRDAP(ctx, pk, arg)
}

func (c CheckRDAPResourceAPI) Delete(ctx context.Context, pk upapi.PrimaryKeyable) error {
	return c.provider.api.Checks().Delete(ctx, pk)
}
