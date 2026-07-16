package provider

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func HeadersSchemaAttribute() schema.Attribute {
	return schema.MapAttribute{
		Description: `A map of HTTP headers where each header name maps to a list of values.
Header names are stored exactly as written; their casing is preserved. Multiple values for the
same header are supported (e.g., { 'Accept': ['application/json', 'text/plain'] }). Defaults to an
empty map if not specified.`,
		ElementType: types.ListType{
			ElemType: types.StringType,
		},
		Optional: true,
		Computed: true,
		Default:  mapdefault.StaticValue(types.MapValueMust(types.ListType{ElemType: types.StringType}, map[string]attr.Value{})),
	}
}

var HeadersType = types.MapType{
	ElemType: types.ListType{ElemType: types.StringType},
}

type HeadersAttributeAdapter struct{}

func (a HeadersAttributeAdapter) HeadersAttributeContext(ctx context.Context, value types.Map) (string, diag.Diagnostics) {
	if !value.Type(ctx).Equal(HeadersType) {
		return "", diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Type expectation failed",
				fmt.Sprintf("expected %s, got %s", HeadersType, value.Type(ctx)),
			),
		}
	}
	if value.IsNull() || value.IsUnknown() {
		return "", nil
	}
	header := make(http.Header)
	diags := value.ElementsAs(ctx, &header, false)
	if diags.HasError() {
		return "", diags
	}
	buf := *new(bytes.Buffer)
	_ = header.Write(&buf)
	return buf.String(), nil
}

func (a HeadersAttributeAdapter) HeadersAttributeValue(in string) (types.Map, error) {
	if strings.TrimSpace(in) == "" {
		return types.MapValueMust(HeadersType.ElementType(), map[string]attr.Value{}), nil
	}
	values := make(map[string][]attr.Value)
	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			return types.MapNull(HeadersType.ElementType()), fmt.Errorf("malformed header line: %q", line)
		}
		name := line[:idx]
		value := strings.TrimSpace(line[idx+1:])
		values[name] = append(values[name], types.StringValue(value))
	}
	if err := scanner.Err(); err != nil {
		return types.MapNull(HeadersType.ElementType()), err
	}
	elems := make(map[string]attr.Value, len(values))
	for name, elem := range values {
		elems[name] = types.ListValueMust(types.StringType, elem)
	}
	return types.MapValueMust(HeadersType.ElementType(), elems), nil
}
