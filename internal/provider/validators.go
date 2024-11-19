package provider

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func OneOfStringValidator(s []string) validator.String {
	return &oneOfStringValidator{
		v: s,
	}
}

type oneOfStringValidator struct {
	v    []string
	m    map[string]struct{}
	once sync.Once
}

func (v *oneOfStringValidator) init() {
	v.m = make(map[string]struct{})
	for _, s := range v.v {
		v.m[s] = struct{}{}
	}
}

func (v *oneOfStringValidator) Description(_ context.Context) string {
	return fmt.Sprintf("one of the following values: %v", v.v)
}

func (v *oneOfStringValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("one of the following values: `%v`", v.v)
}

func (v *oneOfStringValidator) ValidateString(_ context.Context, rq validator.StringRequest, rs *validator.StringResponse) {
	v.once.Do(v.init)
	if rq.ConfigValue.IsNull() {
		return
	}
	if rq.ConfigValue.IsUnknown() {
		return
	}
	if _, ok := v.m[rq.ConfigValue.ValueString()]; !ok {
		rs.Diagnostics.AddAttributeError(
			rq.Path,
			"Invalid value",
			fmt.Sprintf("value must be one of: %v", v.v),
		)
	}
}

func HostnameValidator() validator.String {
	return hostnameValidator{}
}

type hostnameValidator struct {
	zoyaDescriber
}

var hostnameRE = regexp.MustCompile(`^(?:(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}|(?:\d{1,3}\.){3}\d{1,3}|(?:[a-fA-F0-9:]+:+)+[a-fA-F0-9]+)$`)

func (hostnameValidator) ValidateString(_ context.Context, rq validator.StringRequest, rs *validator.StringResponse) {
	if rq.ConfigValue.IsNull() {
		return
	}
	if rq.ConfigValue.IsUnknown() {
		return
	}
	s := rq.ConfigValue.ValueString()
	if !hostnameRE.MatchString(s) {
		rs.Diagnostics.AddAttributeError(
			rq.Path,
			"Invalid value",
			"value must be a valid DNS hostname",
		)
	}
}

func URLValidator() validator.String {
	return urlValidator{}
}

type urlValidator struct {
	zoyaDescriber
}

func (urlValidator) ValidateString(_ context.Context, rq validator.StringRequest, rs *validator.StringResponse) {
	if rq.ConfigValue.IsNull() {
		return
	}
	if rq.ConfigValue.IsUnknown() {
		return
	}
	if u, err := url.Parse(rq.ConfigValue.ValueString()); err != nil {
		rs.Diagnostics.AddAttributeError(
			rq.Path,
			"Invalid value",
			"value must be a valid URL",
		)
	} else if u.Scheme == "" {
		rs.Diagnostics.AddAttributeError(
			rq.Path,
			"Invalid value",
			"value must be a valid URL with a scheme",
		)
	}
}

// PortMatchConfigValidator resource validator checks that address custom port and Port property is same
//
// This validator failes if in the url address property host has explicit port definition
// and Port property in the resource doesn't match it.
func PortMatchConfigValidator(urlP, portP path.Path) *portMatchConfigValidator {
	return &portMatchConfigValidator{
		urlPath:  urlP,
		portPath: portP,
	}
}

type portMatchConfigValidator struct {
	urlPath  path.Path
	portPath path.Path
}

func (p portMatchConfigValidator) Description(context.Context) string {
	return "Port value should match address Host custom port when last is defined"
}

func (p portMatchConfigValidator) MarkdownDescription(ctx context.Context) string {
	return p.Description(ctx)
}

func (p portMatchConfigValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = p.Validate(ctx, req.Config)
}

func (p portMatchConfigValidator) Validate(ctx context.Context, config tfsdk.Config) (diags diag.Diagnostics) {
	// get url address host value and extract port from it if it exists
	var v attr.Value
	getAttributeDiags := config.GetAttribute(ctx, p.urlPath, &v)
	diags.Append(getAttributeDiags...)

	if getAttributeDiags.HasError() || v.IsUnknown() || v.IsNull() {
		return
	}

	urlType, ok := v.(types.String)
	if !ok {
		diags.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			p.urlPath, "path contains non string value", v.String(),
		))
		return
	}

	urlValue, err := url.Parse(urlType.ValueString())
	if err != nil {
		diags.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			p.urlPath, "broken URL", urlType.ValueString(),
		))
		return
	}
	urlPortValue := urlValue.Port()

	// get Port value if it exists
	getAttributeDiags = config.GetAttribute(ctx, p.portPath, &v)
	if getAttributeDiags.HasError() {
		diags.Append(getAttributeDiags...)
		return
	}

	// if Port field is not defined in the resource, but address url host
	// property contains custom port, it is not valid combination
	if (v.IsUnknown() || v.IsNull()) && urlPortValue != "" {
		diags.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			p.portPath, p.Description(ctx),
		))
		return
	}

	portType, ok := v.(types.Int64)
	if !ok {
		diags.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			p.portPath, "path contains non number type", v.String(),
		))
		return
	}

	// No custom ports defined, ok
	if urlPortValue == "" && portType.ValueInt64() == 0 {
		return diags
	}

	// Custom port in the address url port should match with Port property in the resource
	portValue := strconv.FormatInt(portType.ValueInt64(), 10)
	if urlPortValue != portValue {
		diags.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			p.portPath, p.Description(ctx),
		))
	}

	return diags
}
