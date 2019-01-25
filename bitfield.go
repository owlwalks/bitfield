package bitfield

import (
	"log"
	"reflect"
	"strconv"
	"sync"
)

// Upack is unpacking function
type Upack func(in []byte, curr int)

type def struct {
	kind  reflect.Kind
	len   int
	name  string
	uname string
	upack Upack
}

var (
	registry struct {
		index map[string]def
	}
	rMask = [...]uint{0, 128, 192, 224, 240, 248, 252, 254, 0}
	lMask = [...]uint{0, 127, 63, 31, 15, 7, 3, 1, 0}
	once  sync.Once
)

// Register index struct fields for un(packing), should use in init()
func Register(v interface{}) {
	once.Do(func() {
		registry.index = make(map[string]def)
	})

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct && t.Kind() != reflect.Func {
		log.Printf("Register expects struct or func, got %v\n", t.Kind())
		return
	}

	name := t.Name()
	// stop if already defined
	if _, ok := registry.index[name]; ok {
		return
	}
	if t.Kind() == reflect.Func {
		registry.index[name] = def{upack: v.(Upack)}
		return
	}

	num := t.NumField()
	// store struct num fields
	registry.index[name] = def{len: num}
	for i := 0; i < num; i++ {
		registerField(i, name, t.Field(i))
	}
}

func registerField(index int, name string, sf reflect.StructField) {
	sft := sf.Type
	if sft.Kind() == reflect.Struct && sft.Name() != "" {
		if _, ok := registry.index[sf.Name]; !ok {
			log.Printf("%v needs to be registered before %v\n", sf.Name, name)
		}
		return
	}

	switch sft.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Struct: /* Only satisfy struct{} */
		if slen, ok := sf.Tag.Lookup("len"); ok {
			if len, err := strconv.Atoi(slen); err != nil {
				log.Printf("%v has invalid len tag\n", sf.Name)
			} else {
				registry.index[name+strconv.Itoa(index)] = def{kind: sft.Kind(), len: len}
			}
		}
	case reflect.Slice,
		reflect.String:
		if sft.Elem().Kind() == reflect.Uint8 || sft.Kind() == reflect.String {
			if slen, ok := sf.Tag.Lookup("len"); ok {
				if len, err := strconv.Atoi(slen); err != nil {
					log.Printf("%v has invalid len tag\n", sf.Name)
				} else {
					registry.index[name+strconv.Itoa(index)] = def{kind: sft.Kind(), len: len}
					return
				}
			}
		}
		/* len tag on string or []byte slice will overule upack and pack tags */
		fallthrough
	case reflect.Array:
		if uname, ok := sf.Tag.Lookup("upack"); ok {
			if _, existed := registry.index[uname]; !existed {
				log.Printf("Unpack func %v needs to be registered before %v\n", uname, name)
			} else {
				registry.index[name+strconv.Itoa(index)] = def{kind: sft.Kind(), uname: uname}
			}
		}
	default:
		log.Printf("%v (%v) of %v is ignored\n", sf.Name, sft.Kind(), name)
	}
}

func Unpack(dst interface{}, src []byte) {
	v := reflect.Indirect(reflect.ValueOf(dst))
	name := v.Type().Name()
	def, ok := registry.index[name]
	if !ok || v.Kind() != reflect.Struct {
		log.Printf("%v needs to be registered\n", name)
		return
	}

	for fIndex, byteIndex, bitIndex := 0, 0, 0; fIndex < def.len; fIndex++ {
		if fDef, ok := registry.index[name+strconv.Itoa(fIndex)]; ok {
			fVal := v.Field(fIndex)
			val := uint(0)
			l := fDef.len
			switch fDef.kind {
			case reflect.Bool:
				fVal.SetBool(src[byteIndex]&(1<<uint(bitIndex)) != 0)
				byteIndex, bitIndex = nextIndices(byteIndex, bitIndex, fDef.len)
			case reflect.Uint64, reflect.Uint, reflect.Int64, reflect.Int:
				if l > 24 {
					val |= toByte(src, byteIndex, bitIndex, 8) << 24
					byteIndex, bitIndex = nextIndices(byteIndex, bitIndex, 8)
					l -= 8
				}
				fallthrough
			case reflect.Uint32, reflect.Int32:
				if l > 16 {
					val |= toByte(src, byteIndex, bitIndex, 8) << 16
					byteIndex, bitIndex = nextIndices(byteIndex, bitIndex, 8)
					l -= 8
				}
				fallthrough
			case reflect.Uint16, reflect.Int16:
				if l > 8 {
					val |= toByte(src, byteIndex, bitIndex, 8) << 8
					byteIndex, bitIndex = nextIndices(byteIndex, bitIndex, 8)
					l -= 8
				}
				fallthrough
			case reflect.Uint8, reflect.Int8:
				val |= toByte(src, byteIndex, bitIndex, l)
				byteIndex, bitIndex = nextIndices(byteIndex, bitIndex, l)
				if fDef.kind == reflect.Uint8 {
					fVal.SetUint(uint64(val))
				} else {
					fVal.SetInt(int64(val))
				}
			default:
				if fDef.len == 0 && fDef.kind == reflect.Struct {
					byteIndex++
					bitIndex = 0
				} else {
					byteIndex, bitIndex = nextIndices(byteIndex, bitIndex, fDef.len)
				}
			}
		}
	}
}

func toByte(src []byte, byteIndex, bitIndex, bitLen int) uint {
	rShift := 8 - bitIndex - bitLen
	if rShift < 0 {
		rShift = 0
	}

	lShift := bitLen - 8 + bitIndex
	if lShift < 0 || lShift > 7 {
		lShift = 0
	}

	val := (((uint(src[byteIndex]) | rMask[bitIndex]) ^ rMask[bitIndex]) >> uint(rShift)) << uint(lShift)

	if lShift > 0 && (byteIndex+1 <= len(src)-1) {
		val |= (((uint(src[byteIndex+1]) | lMask[lShift]) ^ lMask[lShift]) >> uint(8-lShift))
	}

	return uint(val)
}

func nextIndices(byteIndex, bitIndex, len int) (int, int) {
	bitIndex += len
	return byteIndex + (bitIndex / 8), bitIndex % 8
}
