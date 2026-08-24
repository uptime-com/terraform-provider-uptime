package provider

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

const (
	useIPVersionJSONTag   = "msp_use_ip_version"
	useIPVersionAttribute = "use_ip_version"
)

var useIPVersionResources = []string{
	"uptime_check_api",
	"uptime_check_http",
	"uptime_check_icmp",
	"uptime_check_imap",
	"uptime_check_ntp",
	"uptime_check_pop",
	"uptime_check_smtp",
	"uptime_check_ssh",
	"uptime_check_tcp",
	"uptime_check_udp",
}

func upapiArgumentType(r fwresource.Resource) (reflect.Type, bool) {
	rt := reflect.TypeOf(r)
	if rt == nil || rt.Kind() != reflect.Struct {
		return nil, false
	}
	field, ok := rt.FieldByName("mod")
	if !ok || field.Type.Kind() != reflect.Interface {
		return nil, false
	}
	method, ok := field.Type.MethodByName("ToAPIArgument")
	if !ok || method.Type.NumOut() == 0 {
		return nil, false
	}
	out := method.Type.Out(0)
	if out.Kind() != reflect.Pointer || out.Elem().Kind() != reflect.Struct {
		return nil, false
	}
	return out.Elem(), true
}

func structHasJSONTag(t reflect.Type, tag string) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if name, _, _ := strings.Cut(field.Tag.Get("json"), ","); name == tag {
			return true
		}
		if field.Anonymous && structHasJSONTag(field.Type, tag) {
			return true
		}
	}
	return false
}

func useIPVersionSchemas(t *testing.T) map[string]schema.Attribute {
	t.Helper()
	ctx := context.Background()
	p := &providerImpl{}

	found := map[string]schema.Attribute{}
	for _, newResource := range p.Resources(ctx) {
		r := newResource()

		arg, ok := upapiArgumentType(r)
		if !ok || !structHasJSONTag(arg, useIPVersionJSONTag) {
			continue
		}

		metaResp := fwresource.MetadataResponse{}
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "uptime"}, &metaResp)

		schemaResp := fwresource.SchemaResponse{}
		r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)

		attr, ok := schemaResp.Schema.Attributes[useIPVersionAttribute]
		if !ok {
			t.Errorf("%s: %s accepts %s but the resource schema has no %q attribute - "+
				"users cannot select IPv4 or IPv6 and the value never reaches state",
				metaResp.TypeName, arg.String(), useIPVersionJSONTag, useIPVersionAttribute)
			continue
		}
		found[metaResp.TypeName] = attr
	}
	return found
}

func TestResourcesExposeUseIPVersionWhenTheAPIAcceptsIt(t *testing.T) {
	found := useIPVersionSchemas(t)

	names := make([]string, 0, len(found))
	for name := range found {
		names = append(names, name)
	}
	sort.Strings(names)

	require.Equal(t, useIPVersionResources, names,
		"the set of resources reached by this guard changed - if a resource was added or "+
			"reworked, update useIPVersionResources; a shrinking set means the guard went blind")
}

func TestUseIPVersionRejectsValuesTheAPIDoesNotAccept(t *testing.T) {
	ctx := context.Background()

	schemas := useIPVersionSchemas(t)
	require.NotEmpty(t, schemas,
		"no resource was reached, so every assertion below is vacuous")

	for name, attr := range schemas {
		stringAttr, ok := attr.(schema.StringAttribute)
		if !ok {
			t.Errorf("%s: %q is %T, not a schema.StringAttribute, so it shares neither the "+
				"validator nor the default of the other check resources", name, useIPVersionAttribute, attr)
			continue
		}
		require.NotEmpty(t, stringAttr.Validators, "%s: %q has no validator", name, useIPVersionAttribute)

		require.NotNilf(t, stringAttr.Default, "%s: %q has no default", name, useIPVersionAttribute)
		defaultResp := &defaults.StringResponse{}
		stringAttr.Default.DefaultString(ctx, defaults.StringRequest{}, defaultResp)
		require.Equalf(t, "", defaultResp.PlanValue.ValueString(),
			"%s: %q must default to \"\" (Any); any other default would silently repin checks "+
				"that omit the attribute", name, useIPVersionAttribute)

		for _, value := range []string{"IPV5", "IPv4", "ipv6"} {
			rq := validator.StringRequest{ConfigValue: types.StringValue(value)}
			rs := &validator.StringResponse{}
			for _, v := range stringAttr.Validators {
				v.ValidateString(ctx, rq, rs)
			}
			require.True(t, rs.Diagnostics.HasError(),
				"%s: %q accepted %q, but the API only stores \"\", \"IPV4\" and \"IPV6\"",
				name, useIPVersionAttribute, value)
		}

		for _, value := range []string{"", "IPV4", "IPV6"} {
			rq := validator.StringRequest{ConfigValue: types.StringValue(value)}
			rs := &validator.StringResponse{}
			for _, v := range stringAttr.Validators {
				v.ValidateString(ctx, rq, rs)
			}
			require.False(t, rs.Diagnostics.HasError(),
				"%s: %q rejected %q, which the API accepts", name, useIPVersionAttribute, value)
		}
	}
}

var useIPVersionAdapters = map[string]any{
	"uptime_check_api":  CheckAPIResourceModelAdapter{},
	"uptime_check_http": CheckHTTPResourceModelAdapter{},
	"uptime_check_icmp": CheckICMPResourceModelAdapter{},
	"uptime_check_imap": CheckIMAPResourceModelAdapter{},
	"uptime_check_ntp":  CheckNTPResourceModelAdapter{},
	"uptime_check_pop":  CheckPOPResourceModelAdapter{},
	"uptime_check_smtp": CheckSMTPResourceModelAdapter{},
	"uptime_check_ssh":  CheckSSHResourceModelAdapter{},
	"uptime_check_tcp":  CheckTCPResourceModelAdapter{},
	"uptime_check_udp":  CheckUDPResourceModelAdapter{},
}

func fieldByTag(v reflect.Value, key, name string) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if tag, _, _ := strings.Cut(t.Field(i).Tag.Get(key), ","); tag == name {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}

func TestEveryAdapterSendsUseIPVersionToTheAPI(t *testing.T) {
	require.ElementsMatch(t, useIPVersionResources, keysOf(useIPVersionAdapters))

	for name, adapter := range useIPVersionAdapters {
		for _, value := range []string{"", "IPV4", "IPV6"} {
			method := reflect.ValueOf(adapter).MethodByName("ToAPIArgument")
			model := reflect.New(method.Type().In(0)).Elem()

			field, ok := fieldByTag(model, "tfsdk", useIPVersionAttribute)
			require.Truef(t, ok, "%s: model has no field tagged tfsdk:%q", name, useIPVersionAttribute)
			field.Set(reflect.ValueOf(types.StringValue(value)))

			out := method.Call([]reflect.Value{model})
			require.Truef(t, out[1].IsNil(), "%s: ToAPIArgument returned an error", name)

			sent, ok := fieldByTag(out[0].Elem(), "json", useIPVersionJSONTag)
			require.Truef(t, ok, "%s: argument struct has no field tagged json:%q", name, useIPVersionJSONTag)
			require.Falsef(t, sent.IsNil(),
				"%s: use_ip_version %q never reaches the request body; a nil pointer is dropped by "+
					"omitempty, so the stored value survives a PATCH and the apply fails as inconsistent",
				name, value)
			require.Equalf(t, value, sent.Elem().String(), "%s: wrong value sent", name)
		}
	}
}

func TestEveryAdapterReadsUseIPVersionFromTheAPI(t *testing.T) {
	for name, adapter := range useIPVersionAdapters {
		method := reflect.ValueOf(adapter).MethodByName("FromAPIResult")
		out := method.Call([]reflect.Value{reflect.ValueOf(upapi.Check{UseIPVersion: "IPV6"})})
		require.Truef(t, out[1].IsNil(), "%s: FromAPIResult returned an error", name)

		field, ok := fieldByTag(out[0].Elem(), "tfsdk", useIPVersionAttribute)
		require.Truef(t, ok, "%s: model has no field tagged tfsdk:%q", name, useIPVersionAttribute)

		got, ok := field.Interface().(types.String)
		require.Truef(t, ok, "%s: %q is not a types.String", name, useIPVersionAttribute)
		require.Equalf(t, "IPV6", got.ValueString(),
			"%s: the value the API returned never reaches state, so drift is invisible", name)
	}
}

func keysOf(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
