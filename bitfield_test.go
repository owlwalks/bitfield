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

func TestUnpack(t *testing.T) {
	type s struct {
		_  struct{} `len:"4"`
		F1 int8     `len:"6"`
		_  struct{} `len:"6"`
		F2 int8     `len:"2"`
		_  struct{} `len:"5"`
		F3 int8     `len:"1"`
		F4 bool     `len:"1"`
		F5 int8     `len:"7"`
		F6 int16    `len:"3"`
		F7 int16    `len:"13"`
	}
	s1 := s{}
	Register(s1)
	Unpack(&s1, []byte{0x97, 0x98, 0xD2, 0xB2, 0xCA, 0x28})
	if s1.F1 != 30 {
		t.Errorf("Expect %06b, got %06b\n", 30, s1.F1)
	}
	if s1.F2 != 3 {
		t.Errorf("Expect %02b, got %02b\n", 3, s1.F2)
	}
	if s1.F3 != 0 {
		t.Errorf("Expect %b, got %b\n", 0, s1.F3)
	}
	if s1.F4 {
		t.Errorf("Expect true, got false\n")
	}
	if s1.F5 != 50 {
		t.Errorf("Expect %07b, got %07b\n", 50, s1.F5)
	}
	if s1.F6 != 6 {
		t.Errorf("Expect %03b, got %03b\n", 6, s1.F6)
	}
	if s1.F7 != 20744 {
		t.Errorf("Expect %13b, got %13b\n", 20744, s1.F7)
	}
}
