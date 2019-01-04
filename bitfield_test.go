package bitfield

import (
	"reflect"
	"testing"
)

func Test_registerField(t *testing.T) {
	registry.index = make(map[string]def)
	registerField(0, "struct", reflect.StructField{
		Type: reflect.TypeOf(int(0)),
		Tag:  `len:"4"`,
	})
	if registry.index["struct0"].kind != reflect.Int || registry.index["struct0"].len != 4 {
		t.Error("Unmatched kind and len on struct0")
	}
	registerField(1, "struct", reflect.StructField{
		Type: reflect.TypeOf([]byte(nil)),
		Tag:  `len:"0"`,
	})
	if registry.index["struct1"].kind != reflect.Slice || registry.index["struct1"].len != 0 {
		t.Error("Unmatched kind and len on struct1")
	}
}
