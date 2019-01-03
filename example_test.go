package bitfield_test

import (
	"testing"
)

type dnsmessage struct {
	id        uint16   `len:"16"` // number of bits
	qr        bool     `len:"1"`
	opcode    uint8    `len:"4"`
	aa        bool     `len:"1"`
	tc        bool     `len:"1"`
	rd        bool     `len:"1"`
	ra        bool     `len:"1"`
	_         struct{} `len:"1"` // filling to nearest byte
	ad        bool     `len:"1"`
	cd        bool     `len:"1"`
	rcodes    uint8    `len:"4"`
	qdcount   uint16   `len:"16"`
	ancount   uint16   `len:"16"`
	nscount   uint16   `len:"16"`
	questions []dnsquestion
	answer    []dnsanswer
}

type dnsquestion struct {
	qname  []byte `upack:"nameUpack" pack:"namePack"`
	qtype  uint16 `len:"16"`
	qclass uint16 `len:"16"`
}

type dnsanswer struct {
	name   []byte `upack:"anameUpack"`
	rrtype uint16 `len:"16"`
}

func nameUpack(in []byte, curr int) ([]byte, error) {
	var ret []byte
	for i, delim := curr, 0; i < len(in); i++ {
		if i == curr {
			delim = curr + int(in[i])
			continue
		}
		if in[i] != 0x00 {
			if i == delim+1 {
				ret = append(ret, '.')
				delim = i + int(in[i])
			} else {
				ret = append(ret, in[i])
			}
		} else {
			break
		}
	}
	return ret, nil
}

func anameUpack(in []byte, curr int) ([]byte, error) {
	var ret []byte
	// TODO consider endianness
	for i, delim := curr, 0; i < len(in); i++ {
		if in[i]&0xC0 == 0xC0 {
			// get the pos
			offset := (int(in[i])^0xC0)<<8 | int(in[i+1])
			ptr, _ := nameUpack(in, offset)
			ret = append(ret, '.')
			ret = append(ret, ptr...)
			return ret, nil
		} else {
			if i == curr {
				delim = int(in[i])
				continue
			}
			if in[i] != 0x00 {
				if i == delim+1 {
					ret = append(ret, '.')
					delim = i + int(in[i])
				} else {
					ret = append(ret, in[i])
				}
			}
		}
	}

	return ret, nil
}

func Test_nameUpack(t *testing.T) {
	b := []byte{
		0x01, 0x66, 0x03, 0x69, 0x73, 0x69, 0x04, 0x61,
		0x72, 0x70, 0x61, 0x00, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	u, err := nameUpack(b, 0)
	if string(u) != "f.isi.arpa" {
		t.Errorf("nameUpack() error = %v", err)
	}
}

func Test_anameUpack(t *testing.T) {
	b := []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0x01, 0x66, 0x03, 0x69,
		0x73, 0x69, 0x04, 0x61, 0x72, 0x70, 0x61, 0x00,
		0xFF, 0xFF, 0x03, 0x66, 0x6f, 0x6f, 0xC0, 0x14,
		0xFF, 0xFF, 0xC0, 0x1A, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	u, err := anameUpack(b, 34)
	if string(u) != "foo.f.isi.arpa" {
		t.Errorf("anameUpack() error = %v", err)
	}

	u, err = anameUpack(b, 42)
	if string(u) != ".arpa" {
		t.Errorf("anameUpack() error = %v", err)
	}
}
