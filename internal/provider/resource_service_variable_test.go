package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

// TestServiceVariableSupportsImport pins that the resource is wired up as importable.
// It previously returned a bare APIResource while carrying an ImportState method on the
// API struct - which APIResource holds in a named field, so the method was never promoted
// and terraform import failed with "Resource Import Not Implemented" (SYS-1303).
func TestServiceVariableSupportsImport(t *testing.T) {
	r := NewServiceVariableResource(context.Background(), &providerImpl{})
	if _, ok := r.(fwresource.ResourceWithImportState); !ok {
		t.Fatalf("%T does not implement resource.ResourceWithImportState", r)
	}
}

// importStateForTest runs the resource's import handler against its real schema, so a
// parent attribute that is misnamed or retyped fails in the test rather than at the
// customer's terminal.
func importStateForTest(t *testing.T, api upapi.API, id string) (*fwresource.ImportStateResponse, ServiceVariableResourceModel) {
	t.Helper()
	ctx := context.Background()
	r := NewServiceVariableResource(ctx, &providerImpl{api: api})

	schemaResp := fwresource.SchemaResponse{}
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	resp := &fwresource.ImportStateResponse{State: tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}}

	importer, ok := r.(fwresource.ResourceWithImportState)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithImportState", r)
	}
	importer.ImportState(ctx, fwresource.ImportStateRequest{ID: id}, resp)

	var got ServiceVariableResourceModel
	if !resp.Diagnostics.HasError() {
		if diags := resp.State.Get(ctx, &got); diags.HasError() {
			t.Fatalf("reading imported state: %v", diags.Errors())
		}
	}
	return resp, got
}

// TestServiceVariableImportState covers the composite import: service_id must land in
// state because the API does not return it (Read sources it from the model, and it is
// Required + RequiresReplace, so an unset one makes the next plan propose a replacement),
// and a service_id naming some other check must be rejected rather than silently accepted
// - nothing server-side ties the two halves together (SYS-1303).
func TestServiceVariableImportState(t *testing.T) {
	api := stubServiceVariablesAPI{
		get:   &upapi.ServiceVariable{ID: 41218, Service: "arya-sanity-euw"},
		check: &upapi.Check{PK: 5891524, Name: "arya-sanity-euw"},
	}

	t.Run("writes both halves into state", func(t *testing.T) {
		resp, got := importStateForTest(t, api, "5891524:41218")
		if resp.Diagnostics.HasError() {
			t.Fatalf("import failed: %v", resp.Diagnostics.Errors())
		}
		if got.ServiceID.ValueInt64() != 5891524 {
			t.Errorf("service_id = %d, want 5891524", got.ServiceID.ValueInt64())
		}
		if got.ID.ValueInt64() != 41218 {
			t.Errorf("id = %d, want 41218", got.ID.ValueInt64())
		}
	})

	t.Run("rejects a service_id owned by another check", func(t *testing.T) {
		mismatched := api
		mismatched.check = &upapi.Check{PK: 999, Name: "some-other-check"}
		resp, _ := importStateForTest(t, mismatched, "999:41218")
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected an error: service_id names a check that does not own the variable")
		}
		if summary := resp.Diagnostics.Errors()[0].Summary(); summary != "Import ID Mismatch" {
			t.Errorf("summary = %q, want \"Import ID Mismatch\"", summary)
		}
	})

	t.Run("imports when the API omits either name", func(t *testing.T) {
		blank := api
		blank.get = &upapi.ServiceVariable{ID: 41218}
		resp, got := importStateForTest(t, blank, "5891524:41218")
		if resp.Diagnostics.HasError() {
			t.Fatalf("a blank name must not block import: %v", resp.Diagnostics.Errors())
		}
		if got.ServiceID.ValueInt64() != 5891524 {
			t.Errorf("service_id = %d, want 5891524", got.ServiceID.ValueInt64())
		}
	})
}

// TestServiceVariablePreservePlanValues verifies that Required attributes the update
// endpoint may omit (credential_id, property_name, variable_name) are restored from the
// plan when the API result comes back empty, preventing "Provider produced inconsistent
// result after apply".
func TestServiceVariablePreservePlanValues(t *testing.T) {
	adapter := ServiceVariableResourceModelAdapter{}

	plan := &ServiceVariableResourceModel{
		CredentialID: types.Int64Value(28966),
		PropertyName: types.StringValue("secret"),
		VariableName: types.StringValue("token_raven_token"),
	}

	t.Run("restores omitted values", func(t *testing.T) {
		result := &ServiceVariableResourceModel{
			CredentialID: types.Int64Value(0),
			PropertyName: types.StringValue(""),
			VariableName: types.StringValue(""),
		}
		got := adapter.PreservePlanValues(result, plan)
		if got.CredentialID.ValueInt64() != 28966 {
			t.Errorf("credential_id: got %d, want 28966", got.CredentialID.ValueInt64())
		}
		if got.PropertyName.ValueString() != "secret" {
			t.Errorf("property_name: got %q, want \"secret\"", got.PropertyName.ValueString())
		}
		if got.VariableName.ValueString() != "token_raven_token" {
			t.Errorf("variable_name: got %q, want \"token_raven_token\"", got.VariableName.ValueString())
		}
	})

	t.Run("keeps non-empty result values", func(t *testing.T) {
		result := &ServiceVariableResourceModel{
			CredentialID: types.Int64Value(28965),
			PropertyName: types.StringValue("password"),
			VariableName: types.StringValue("other_name"),
		}
		got := adapter.PreservePlanValues(result, plan)
		if got.CredentialID.ValueInt64() != 28965 {
			t.Errorf("credential_id: got %d, want 28965", got.CredentialID.ValueInt64())
		}
		if got.PropertyName.ValueString() != "password" {
			t.Errorf("property_name: got %q, want \"password\"", got.PropertyName.ValueString())
		}
		if got.VariableName.ValueString() != "other_name" {
			t.Errorf("variable_name: got %q, want \"other_name\"", got.VariableName.ValueString())
		}
	})
}

// TestServiceVariablePreserveReadValues verifies that refresh trusts the API and does
// not backfill from prior state, so a UI-side removal of the credential link (which the
// API reports as empty) surfaces as drift instead of being masked (SYS-1284).
func TestServiceVariablePreserveReadValues(t *testing.T) {
	adapter := ServiceVariableResourceModelAdapter{}

	state := &ServiceVariableResourceModel{
		CredentialID: types.Int64Value(28966),
		PropertyName: types.StringValue("secret"),
		VariableName: types.StringValue("token_raven_token"),
	}
	result := &ServiceVariableResourceModel{
		CredentialID: types.Int64Value(0),
		PropertyName: types.StringValue(""),
		VariableName: types.StringValue(""),
	}

	got := adapter.PreserveReadValues(result, state)
	if got.CredentialID.ValueInt64() != 0 {
		t.Errorf("credential_id: got %d, want 0 (must not backfill from state)", got.CredentialID.ValueInt64())
	}
	if got.PropertyName.ValueString() != "" {
		t.Errorf("property_name: got %q, want \"\" (must not backfill from state)", got.PropertyName.ValueString())
	}
	if got.VariableName.ValueString() != "" {
		t.Errorf("variable_name: got %q, want \"\" (must not backfill from state)", got.VariableName.ValueString())
	}
}

// stubServiceVariablesAPI embeds upapi.API and overrides only ServiceVariables so the
// CRUD paths can be exercised without a live client. Get returns get; Create and Update
// return write, or err when it is set. Every other method is left unimplemented by the
// embedded interface and panics if called.
//
// The real client never answers (nil, nil) - endpointCreatorImpl always returns either an
// error or a record - so a test that leaves get or write nil would panic inside provider
// code rather than here. The accessors reject that up front to keep the panic local.
type stubServiceVariablesAPI struct {
	upapi.API
	get   *upapi.ServiceVariable
	write *upapi.ServiceVariable
	check *upapi.Check
	err   error
}

func (s stubServiceVariablesAPI) ServiceVariables() upapi.ServiceVariablesEndpoint {
	return stubServiceVariablesEndpoint{get: s.get, write: s.write, err: s.err}
}

// Checks is reached only by the import handler, which cross-checks that service_id names
// the check the variable reports as its owner.
func (s stubServiceVariablesAPI) Checks() upapi.ChecksEndpoint {
	return stubChecksEndpoint{get: s.check}
}

type stubChecksEndpoint struct {
	upapi.ChecksEndpoint
	get *upapi.Check
}

func (s stubChecksEndpoint) Get(context.Context, upapi.PrimaryKeyable) (*upapi.Check, error) {
	if s.get == nil {
		panic("stubChecksEndpoint: Get called but check is not set")
	}
	return s.get, nil
}

type stubServiceVariablesEndpoint struct {
	upapi.ServiceVariablesEndpoint
	get   *upapi.ServiceVariable
	write *upapi.ServiceVariable
	err   error
}

func (s stubServiceVariablesEndpoint) Get(context.Context, upapi.PrimaryKeyable) (*upapi.ServiceVariable, error) {
	if s.get == nil {
		panic("stubServiceVariablesEndpoint: Get called but get is not set")
	}
	return s.get, nil
}

func (s stubServiceVariablesEndpoint) Create(
	context.Context, upapi.ServiceVariableCreateRequest,
) (*upapi.ServiceVariable, error) {
	return s.writeResult("Create")
}

func (s stubServiceVariablesEndpoint) Update(
	context.Context, upapi.PrimaryKeyable, upapi.ServiceVariableUpdateRequest,
) (*upapi.ServiceVariable, error) {
	return s.writeResult("Update")
}

func (s stubServiceVariablesEndpoint) writeResult(op string) (*upapi.ServiceVariable, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.write == nil {
		panic("stubServiceVariablesEndpoint: " + op + " called but write is not set")
	}
	return s.write, nil
}

// TestServiceVariableReadDeletedDrift verifies that a link the API reports as
// soft-deleted (deleted_at set) is treated as gone on refresh, so the deletion
// surfaces as drift instead of being masked by the backfill (SYS-1284).
func TestServiceVariableReadDeletedDrift(t *testing.T) {
	deletedAt := time.Unix(1_700_000_000, 0)
	api := ServiceVariableResourceAPI{provider: &providerImpl{
		api: stubServiceVariablesAPI{get: &upapi.ServiceVariable{ID: 42, DeletedAt: &deletedAt}},
	}}

	_, err := api.Read(context.Background(), ServiceVariableResourceModel{ID: types.Int64Value(42)})
	if !errors.Is(err, errResourceGone) {
		t.Fatalf("Read of deleted link: got err %v, want errResourceGone", err)
	}
	if !isNotFoundError(err) {
		t.Errorf("isNotFoundError(%v) = false, want true so the resource is dropped from state", err)
	}
}

// TestServiceVariableReadHardDeletedGone verifies that a link removed from the check
// in the UI is treated as gone on refresh. That removal is a hard delete, so the API
// answers a subsequent Get with HTTP 200 and an empty body (no deleted_at), which the
// client decodes as a zero-valued record. Read must surface that as errResourceGone so
// the resource is dropped from state and the next apply recreates the link, instead of
// issuing a no-op update against id 0 that never restores it (SYS-1284 follow-up).
func TestServiceVariableReadHardDeletedGone(t *testing.T) {
	api := ServiceVariableResourceAPI{provider: &providerImpl{
		api: stubServiceVariablesAPI{get: &upapi.ServiceVariable{ID: 0}},
	}}

	_, err := api.Read(context.Background(), ServiceVariableResourceModel{ID: types.Int64Value(42)})
	if !errors.Is(err, errResourceGone) {
		t.Fatalf("Read of hard-deleted link: got err %v, want errResourceGone", err)
	}
	if !isNotFoundError(err) {
		t.Errorf("isNotFoundError(%v) = false, want true so the resource is dropped from state", err)
	}
}

// TestServiceVariableReadLive verifies that a live link is returned as-is and, when
// credential_id is only present in the nested credential object, is recovered from it.
func TestServiceVariableReadLive(t *testing.T) {
	api := ServiceVariableResourceAPI{provider: &providerImpl{
		api: stubServiceVariablesAPI{get: &upapi.ServiceVariable{
			ID:           42,
			Credential:   &upapi.ServiceVariableCredential{ID: 99},
			PropertyName: "password",
			VariableName: "api_password",
		}},
	}}

	got, err := api.Read(context.Background(), ServiceVariableResourceModel{
		ID:        types.Int64Value(42),
		ServiceID: types.Int64Value(7),
	})
	if err != nil {
		t.Fatalf("Read of live link: unexpected error %v", err)
	}
	if got.CredentialID != 99 {
		t.Errorf("credential_id: got %d, want 99 (recovered from nested credential)", got.CredentialID)
	}
	if got.ServiceID != 7 {
		t.Errorf("service_id: got %d, want 7 (preserved from prior state)", got.ServiceID)
	}
}

// TestServiceVariableWriteMissingID verifies that a write whose response carries no ID is
// reported as an error instead of being persisted. The endpoint answers a rejected write
// with HTTP 200 and an empty results object; because id is Computed, Terraform would accept
// the resulting id 0 as a valid apply result and the failure would only surface later, as a
// refresh that reads the resource as gone and recreates it on every plan.
func TestServiceVariableWriteMissingID(t *testing.T) {
	api := ServiceVariableResourceAPI{provider: &providerImpl{
		api: stubServiceVariablesAPI{write: &upapi.ServiceVariable{ID: 0}},
	}}
	arg := ServiceVariableWrapper{
		ServiceID: 7,
		ServiceVariable: upapi.ServiceVariable{
			CredentialID: 99,
			PropertyName: "secret",
			VariableName: "token_arya_sanity",
		},
	}

	t.Run("create", func(t *testing.T) {
		got, err := api.Create(context.Background(), arg)
		if !errors.Is(err, errServiceVariableNoID) {
			t.Fatalf("Create with empty result: got err %v, want errServiceVariableNoID", err)
		}
		if got != nil {
			t.Errorf("Create with empty result: got %+v, want nil so id 0 is never persisted", got)
		}
		if isNotFoundError(err) {
			t.Error("isNotFoundError = true, want false: a failed write must surface as an " +
				"apply error, not silently drop the resource from state")
		}
	})

	t.Run("update", func(t *testing.T) {
		got, err := api.Update(context.Background(), ServiceVariableResourceModel{ID: types.Int64Value(42)}, arg)
		if !errors.Is(err, errServiceVariableNoID) {
			t.Fatalf("Update with empty result: got err %v, want errServiceVariableNoID", err)
		}
		if got != nil {
			t.Errorf("Update with empty result: got %+v, want nil so id 0 is never persisted", got)
		}
		if isNotFoundError(err) {
			t.Error("isNotFoundError = true, want false: a failed write must surface as an " +
				"apply error, not silently drop the resource from state")
		}
	})
}

// TestServiceVariableWriteLive verifies the guard does not reject a normal write, that
// credential_id is still recovered from the nested credential object, and that service_id
// is carried over from the request - the API does not echo it back.
func TestServiceVariableWriteLive(t *testing.T) {
	api := ServiceVariableResourceAPI{provider: &providerImpl{
		api: stubServiceVariablesAPI{write: &upapi.ServiceVariable{
			ID:           42,
			Credential:   &upapi.ServiceVariableCredential{ID: 99},
			PropertyName: "secret",
			VariableName: "token_arya_sanity",
		}},
	}}
	arg := ServiceVariableWrapper{ServiceID: 7}

	assert := func(t *testing.T, got *ServiceVariableWrapper, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("write of live link: unexpected error %v", err)
		}
		if got.ID != 42 {
			t.Errorf("id: got %d, want 42", got.ID)
		}
		if got.CredentialID != 99 {
			t.Errorf("credential_id: got %d, want 99 (recovered from nested credential)", got.CredentialID)
		}
		if got.ServiceID != 7 {
			t.Errorf("service_id: got %d, want 7 (preserved from the request)", got.ServiceID)
		}
	}

	t.Run("create", func(t *testing.T) {
		got, err := api.Create(context.Background(), arg)
		assert(t, got, err)
	})

	t.Run("update", func(t *testing.T) {
		got, err := api.Update(context.Background(), ServiceVariableResourceModel{ID: types.Int64Value(42)}, arg)
		assert(t, got, err)
	})
}

// TestServiceVariableWriteError verifies that a genuine client error is returned as-is and
// not masked by the empty-result guard, which would misreport the cause of the failure.
func TestServiceVariableWriteError(t *testing.T) {
	apiErr := errors.New("connection refused")
	api := ServiceVariableResourceAPI{provider: &providerImpl{
		api: stubServiceVariablesAPI{err: apiErr},
	}}

	t.Run("create", func(t *testing.T) {
		_, err := api.Create(context.Background(), ServiceVariableWrapper{ServiceID: 7})
		if !errors.Is(err, apiErr) {
			t.Fatalf("Create: got err %v, want the client error", err)
		}
	})

	t.Run("update", func(t *testing.T) {
		_, err := api.Update(context.Background(), ServiceVariableResourceModel{ID: types.Int64Value(42)},
			ServiceVariableWrapper{ServiceID: 7})
		if !errors.Is(err, apiErr) {
			t.Fatalf("Update: got err %v, want the client error", err)
		}
	})
}

func TestAccServiceVariableResource(t *testing.T) {
	credentialName := petname.Generate(3, "-")
	password := petname.Generate(1, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_basic"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
				"variable_name":   config.StringVariable("api_password"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "api_password"),
				resource.TestCheckResourceAttr("uptime_service_variable.test", "property_name", "password"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "service_id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "credential_id"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_basic"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
				"variable_name":   config.StringVariable("api_key"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "api_key"),
				resource.TestCheckResourceAttr("uptime_service_variable.test", "property_name", "password"),
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.test", "id",
				),
			),
		},
		{
			// ImportStateVerify is the only check that the composite import really
			// reconstructs full state: service_id comes from the import ID, everything
			// else has to be repopulated by the post-import Read (SYS-1303).
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_basic"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
				"variable_name":   config.StringVariable("api_key"),
			},
			ResourceName:      "uptime_service_variable.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				rs := s.RootModule().Resources["uptime_service_variable.test"]
				if rs == nil {
					return "", fmt.Errorf("uptime_service_variable.test not found in state")
				}
				return fmt.Sprintf("%s:%s",
					rs.Primary.Attributes["service_id"], rs.Primary.Attributes["id"]), nil
			},
		},
	}))
}

// TestAccServiceVariableResource_TokenSwap reproduces the customer scenario: a TOKEN
// credential service variable whose credential_id changes while variable_name stays
// constant.
func TestAccServiceVariableResource_TokenSwap(t *testing.T) {
	credentialNameA := petname.Generate(3, "-")
	credentialNameB := petname.Generate(3, "-")
	token := petname.Generate(2, "-")

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_token_swap"),
			ConfigVariables: config.Variables{
				"credential_name_a": config.StringVariable(credentialNameA),
				"credential_name_b": config.StringVariable(credentialNameB),
				"token":             config.StringVariable(token),
				"use_b":             config.BoolVariable(false),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "token_raven_token"),
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.a", "id",
				),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/_token_swap"),
			ConfigVariables: config.Variables{
				"credential_name_a": config.StringVariable(credentialNameA),
				"credential_name_b": config.StringVariable(credentialNameB),
				"token":             config.StringVariable(token),
				"use_b":             config.BoolVariable(true),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", "token_raven_token"),
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.b", "id",
				),
			),
		},
	}))
}

func TestAccServiceVariableResource_WithDataSource(t *testing.T) {
	credentialName := petname.Generate(3, "-")
	password := petname.Generate(1, "-")
	variableName := "api_password"

	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/datasource_step1"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", credentialName),
				resource.TestCheckResourceAttr("uptime_credential.test", "credential_type", "BASIC"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_service_variable/datasource_step2"),
			ConfigVariables: config.Variables{
				"credential_name": config.StringVariable(credentialName),
				"password":        config.StringVariable(password),
				"variable_name":   config.StringVariable(variableName),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				// Check credential exists
				resource.TestCheckResourceAttr("uptime_credential.test", "display_name", credentialName),
				// Check datasource works
				resource.TestCheckResourceAttrSet("data.uptime_credentials.all", "credentials.#"),
				// Check service variable is created
				resource.TestCheckResourceAttr("uptime_service_variable.test", "variable_name", variableName),
				resource.TestCheckResourceAttr("uptime_service_variable.test", "property_name", "password"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "service_id"),
				resource.TestCheckResourceAttrSet("uptime_service_variable.test", "credential_id"),
				// Verify the credential_id matches the one from datasource
				resource.TestCheckResourceAttrPair(
					"uptime_service_variable.test", "credential_id",
					"uptime_credential.test", "id",
				),
			),
		},
	}))
}
