package provider

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

var hostnameRE = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])`)

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
