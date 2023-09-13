package provider

import (
	"github.com/uptime-com/terraform-provider-uptime/internal/reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

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
