package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

var _ resource.Resource = (*UserResource)(nil)

func NewUserResource(_ context.Context, p *providerImpl) resource.Resource {
	return &UserResource{
		provider: p,
	}
}

type UserResource struct {
	provider *providerImpl
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Account user resource. Import using the user ID: `terraform import uptime_user.example 123`",
		Attributes: map[string]schema.Attribute{
			"id":  IDSchemaAttribute(),
			"url": URLSchemaAttribute(),
			"first_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"last_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_active": schema.BoolAttribute{
				Computed: true,
			},
			"is_primary": schema.BoolAttribute{
				Computed: true,
			},
			"access_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"is_api_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"notify_paid_invoices": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"assigned_subaccounts": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"require_two_factor": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"must_two_factor": schema.BoolAttribute{
				Computed: true,
			},
			"timezone": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type UserResourceModel struct {
	ID                  types.Int64  `tfsdk:"id"`
	URL                 types.String `tfsdk:"url"`
	FirstName           types.String `tfsdk:"first_name"`
	LastName            types.String `tfsdk:"last_name"`
	Email               types.String `tfsdk:"email"`
	Password            types.String `tfsdk:"password"`
	IsActive            types.Bool   `tfsdk:"is_active"`
	IsPrimary           types.Bool   `tfsdk:"is_primary"`
	AccessLevel         types.String `tfsdk:"access_level"`
	IsAPIEnabled        types.Bool   `tfsdk:"is_api_enabled"`
	NotifyPaidInvoices  types.Bool   `tfsdk:"notify_paid_invoices"`
	AssignedSubaccounts types.Set    `tfsdk:"assigned_subaccounts"`
	RequireTwoFactor    types.String `tfsdk:"require_two_factor"`
	MustTwoFactor       types.Bool   `tfsdk:"must_two_factor"`
	Timezone            types.String `tfsdk:"timezone"`
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setAttrAdapter := SetAttributeAdapter[string]{}
	isAPIEnabled := plan.IsAPIEnabled.ValueBool()
	notifyPaidInvoices := plan.NotifyPaidInvoices.ValueBool()

	createReq := upapi.UserCreateRequest{
		FirstName:           plan.FirstName.ValueString(),
		LastName:            plan.LastName.ValueString(),
		Email:               plan.Email.ValueString(),
		Password:            plan.Password.ValueString(),
		AccessLevel:         plan.AccessLevel.ValueString(),
		IsAPIEnabled:        isAPIEnabled,
		NotifyPaidInvoices:  notifyPaidInvoices,
		AssignedSubaccounts: setAttrAdapter.Slice(plan.AssignedSubaccounts),
		RequireTwoFactor:    plan.RequireTwoFactor.ValueString(),
	}

	user, err := r.provider.api.Users().Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create user", err.Error())
		return
	}

	// Map response to state, preserving password from plan
	r.mapUserToModel(user, &plan, plan.Password)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.provider.api.Users().Get(ctx, upapi.PrimaryKey(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user", err.Error())
		return
	}

	// Map response to state, preserving password from current state
	r.mapUserToModel(user, &state, state.Password)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setAttrAdapter := SetAttributeAdapter[string]{}
	isAPIEnabled := plan.IsAPIEnabled.ValueBool()
	notifyPaidInvoices := plan.NotifyPaidInvoices.ValueBool()

	updateReq := upapi.UserUpdateRequest{
		FirstName:           plan.FirstName.ValueString(),
		LastName:            plan.LastName.ValueString(),
		Email:               plan.Email.ValueString(),
		Password:            plan.Password.ValueString(),
		AccessLevel:         plan.AccessLevel.ValueString(),
		IsAPIEnabled:        &isAPIEnabled,
		NotifyPaidInvoices:  &notifyPaidInvoices,
		AssignedSubaccounts: setAttrAdapter.Slice(plan.AssignedSubaccounts),
		RequireTwoFactor:    plan.RequireTwoFactor.ValueString(),
	}

	user, err := r.provider.api.Users().Update(ctx, upapi.PrimaryKey(state.ID.ValueInt64()), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update user", err.Error())
		return
	}

	// Map response to state, preserving password from plan
	r.mapUserToModel(user, &plan, plan.Password)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.provider.api.Users().Delete(ctx, upapi.PrimaryKey(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete user", err.Error())
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStateSimpleID(ctx, req, resp)
}

func (r *UserResource) mapUserToModel(user *upapi.User, model *UserResourceModel, password types.String) {
	setAttrAdapter := SetAttributeAdapter[string]{}

	// Helper to preserve existing value if API returns empty
	stringValueOrKeep := func(apiValue string, currentValue types.String) types.String {
		if apiValue == "" {
			if currentValue.IsNull() || currentValue.IsUnknown() {
				return types.StringNull()
			}
			// Preserve the current value if API returns empty
			return currentValue
		}
		return types.StringValue(apiValue)
	}

	model.ID = types.Int64Value(user.PK)
	model.URL = types.StringValue(user.URL)
	model.FirstName = stringValueOrKeep(user.FirstName, model.FirstName)
	model.LastName = stringValueOrKeep(user.LastName, model.LastName)
	// For email, always use API value since it's required and should never be empty
	if user.Email != "" {
		model.Email = types.StringValue(user.Email)
	}
	model.Password = password // Preserve password from plan/state
	model.IsActive = types.BoolValue(user.IsActive)
	model.IsPrimary = types.BoolValue(user.IsPrimary)
	model.AccessLevel = stringValueOrKeep(user.AccessLevel, model.AccessLevel)
	model.IsAPIEnabled = types.BoolValue(user.IsAPIEnabled)
	model.NotifyPaidInvoices = types.BoolValue(user.NotifyPaidInvoices)
	model.AssignedSubaccounts = setAttrAdapter.SliceValue(user.AssignedSubaccounts)
	model.RequireTwoFactor = stringValueOrKeep(user.RequireTwoFactor, model.RequireTwoFactor)
	model.MustTwoFactor = types.BoolValue(user.MustTwoFactor)
	model.Timezone = stringValueOrKeep(user.Timezone, model.Timezone)
}
