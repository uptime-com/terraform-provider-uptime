package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func NewPushNotificationResource(_ context.Context, p *providerImpl) resource.Resource {
	return &pushNotificationResource{
		provider: p,
	}
}

type pushNotificationResource struct {
	provider *providerImpl
}

func (r *pushNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_notification"
}

func (r *pushNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Push notification profile resource for mobile device notifications",
		Attributes: map[string]schema.Attribute{
			"id":  IDSchemaAttribute(),
			"url": URLSchemaAttribute(),
			"app_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Authentication key from the mobile app",
			},
			"uuid": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the mobile device",
			},
			"device_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the mobile device",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "Display name for the push notification profile",
			},
			"user": schema.StringAttribute{
				Computed:    true,
				Description: "User associated with the push notification profile",
			},
			"contact_groups": ContactGroupsSchemaAttribute(),
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the profile was created",
			},
			"modified_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the profile was last modified",
			},
		},
	}
}

type PushNotificationResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	AppKey        types.String `tfsdk:"app_key"`
	UUID          types.String `tfsdk:"uuid"`
	DeviceName    types.String `tfsdk:"device_name"`
	DisplayName   types.String `tfsdk:"display_name"`
	User          types.String `tfsdk:"user"`
	ContactGroups types.Set    `tfsdk:"contact_groups"`
	CreatedAt     types.String `tfsdk:"created_at"`
	ModifiedAt    types.String `tfsdk:"modified_at"`
}

func (m PushNotificationResourceModel) PrimaryKey() upapi.PrimaryKey {
	return upapi.PrimaryKey(m.ID.ValueInt64())
}

type pushNotificationResourceModelAdapter struct {
	ContactGroupsAttributeAdapter
}

func (a pushNotificationResourceModelAdapter) contactGroupsFromAPI(cg []string) types.Set {
	return a.ContactGroupsValue(cg)
}

func (a pushNotificationResourceModelAdapter) contactGroupsToAPI(set types.Set) []string {
	return a.ContactGroups(set)
}

func (r *pushNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PushNotificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adapter := pushNotificationResourceModelAdapter{}
	createReq := upapi.PushNotificationProfileCreateRequest{
		AppKey:        plan.AppKey.ValueString(),
		UUID:          plan.UUID.ValueString(),
		DeviceName:    plan.DeviceName.ValueString(),
		ContactGroups: adapter.contactGroupsToAPI(plan.ContactGroups),
	}

	profile, err := r.provider.api.PushNotifications().Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create push notification profile", err.Error())
		return
	}

	// Update state with API response, preserving AppKey from plan
	state := PushNotificationResourceModel{
		ID:            types.Int64Value(profile.PK),
		URL:           types.StringValue(profile.URL),
		AppKey:        plan.AppKey, // Preserve from plan
		UUID:          types.StringValue(profile.UUID),
		DeviceName:    types.StringValue(profile.DeviceName),
		DisplayName:   types.StringValue(profile.DisplayName),
		User:          types.StringValue(profile.User),
		ContactGroups: adapter.contactGroupsFromAPI(profile.ContactGroups),
		CreatedAt:     types.StringValue(profile.CreatedAt),
		ModifiedAt:    types.StringValue(profile.ModifiedAt),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *pushNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PushNotificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.provider.api.PushNotifications().Get(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read push notification profile", err.Error())
		return
	}

	adapter := pushNotificationResourceModelAdapter{}
	// Update state with API response, preserving AppKey from current state
	state.ID = types.Int64Value(profile.PK)
	state.URL = types.StringValue(profile.URL)
	// state.AppKey preserved from current state
	state.UUID = types.StringValue(profile.UUID)
	state.DeviceName = types.StringValue(profile.DeviceName)
	state.DisplayName = types.StringValue(profile.DisplayName)
	state.User = types.StringValue(profile.User)
	state.ContactGroups = adapter.contactGroupsFromAPI(profile.ContactGroups)
	state.CreatedAt = types.StringValue(profile.CreatedAt)
	state.ModifiedAt = types.StringValue(profile.ModifiedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *pushNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PushNotificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PushNotificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adapter := pushNotificationResourceModelAdapter{}
	updateReq := upapi.PushNotificationProfileUpdateRequest{
		DeviceName:    plan.DeviceName.ValueString(),
		ContactGroups: adapter.contactGroupsToAPI(plan.ContactGroups),
	}

	profile, err := r.provider.api.PushNotifications().Update(ctx, state, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update push notification profile", err.Error())
		return
	}

	// Update state with API response, preserving AppKey from plan
	state = PushNotificationResourceModel{
		ID:            types.Int64Value(profile.PK),
		URL:           types.StringValue(profile.URL),
		AppKey:        plan.AppKey, // Preserve from plan
		UUID:          types.StringValue(profile.UUID),
		DeviceName:    types.StringValue(profile.DeviceName),
		DisplayName:   types.StringValue(profile.DisplayName),
		User:          types.StringValue(profile.User),
		ContactGroups: adapter.contactGroupsFromAPI(profile.ContactGroups),
		CreatedAt:     types.StringValue(profile.CreatedAt),
		ModifiedAt:    types.StringValue(profile.ModifiedAt),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *pushNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PushNotificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.provider.api.PushNotifications().Delete(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete push notification profile", err.Error())
		return
	}
}
