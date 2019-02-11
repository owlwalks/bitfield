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
		_   struct{} `len:"4"`
		F1  int8     `len:"6"`
		_   struct{} `len:"6"`
		F2  int8     `len:"2"`
		_   struct{} `len:"5"`
		F3  int8     `len:"1"`
		F4  bool     `len:"1"`
		F5  int8     `len:"7"`
		F6  int16    `len:"3"`
		_   struct{} `len:"0"`
		F7  int16    `len:"7"`
		_   struct{} `len:"0"`
		F8  int32    `len:"16"`
		F9  int      `len:"25"`
		F10 int      `len:"33"`
	}
	s1 := s{}
	Register(s1)
	BigEndian.Unpack(&s1, []byte{0x97, 0x98, 0xD2, 0xB2, 0xCA, 0x28, 0x99, 0x3F, 0xD5, 0xE7, 0x70, 0xFC, 0x35, 0xEE, 0x54, 0x58, 0xC6, 0x14})
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
	if s1.F7 != 20 {
		t.Errorf("Expect %013b, got %013b\n", 20, s1.F7)
	}
	if s1.F8 != 39231 {
		t.Errorf("Expect %016b, got %016b\n", 39231, s1.F8)
	}
	if s1.F9 != 3588714497 {
		t.Errorf("Expect %025b, got %025b\n", 3588714497, s1.F9)
	}
	if s1.F10 != 1066961512449 {
		t.Errorf("Expect %033b, got %033b\n", 1066961512449, s1.F10)
	}
}
