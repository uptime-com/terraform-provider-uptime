package provider

import (
	"bytes"
	"encoding/json"
	"github.com/uptime-com/terraform-provider-uptime/internal/reflect"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func ptr[T any](v T) *T {
	return &v
}

func mustErr[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func mustDiag[T any](x T, diags diag.Diagnostics) T {
	if diags.HasError() {
		panic(diags)
	}
	return x
}

func prettyJSON(r io.Reader) io.Reader {
	data := make(map[string]any)
	_ = json.NewDecoder(r).Decode(&data)
	return bytes.NewReader(mustErr(json.MarshalIndent(&data, "", "\t")))
}

func prettyResponse(status string, body io.Reader) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(status)
	buf.WriteRune('\n')
	_, _ = buf.ReadFrom(prettyJSON(body))
	return buf.String()
}

func valueFromAPI(dst, src any) diag.Diagnostics {
	err := reflect.CopyIn(dst, src)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("reflect.CopyIn", err.Error()),
		}
	}
	return nil
}

func valueToAPI(dst, src any) diag.Diagnostics {
	err := reflect.CopyOut(dst, src)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("reflect.CopyOut", err.Error()),
		}
	}
	return nil
}