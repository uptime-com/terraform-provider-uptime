package provider

import (
	"bytes"
	"encoding/json"
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
