package uptime

import (
	"reflect"
	"testing"
)

func TestHeaderConversion(t *testing.T) {
	hmIface := make(map[string]interface{})
	hmIface["Foo"] = "bar"
	hmIface["Baz"] = "bat"

	hm := map[string]string{
		"Foo": "bar",
		"Baz": "bat",
	}

	hs := "Foo: bar\nBaz: bat\n"

	if s := headersMapToString(hmIface); s != hs {
		t.Errorf("headersMaptoString returned '%s', expected '%s'", s, hs)
	}

	if m := headersStringToMap(hs); !reflect.DeepEqual(m, hm) {
		t.Errorf("headersStringToMap returned '%+v', expected '%+v'", m, hm)
	}
}
