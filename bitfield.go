package bitfield

import (
	"log"
	"reflect"
	"strconv"
	"sync"
)

var once sync.Once

// Unpack is unpacking function
type Unpack func(in []byte, curr int)

type def struct {
	kind  reflect.Kind
	len   int
	name  string
	uname string
	upack Unpack
}

var registry struct {
	sync.RWMutex
	index map[string]def
}

// Register index struct fields for un(packing), should use in init()
func Register(v interface{}) {
	once.Do(func() {
		registry.Lock()
		defer registry.Unlock()
		registry.index = make(map[string]def)
	})

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct || t.Kind() != reflect.Func {
		log.Printf("Register expects struct or func, got %v\n", t.Kind())
		return
	}

	name := t.Name()
	registry.Lock()
	defer registry.Unlock()
	// stop if already defined
	if _, ok := registry.index[name]; ok {
		return
	}
	if t.Kind() == reflect.Func {
		registry.index[name] = def{upack: v.(Unpack)}
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
	// reflect.Float32,
	// reflect.Float64,
	// reflect.Complex64,
	// reflect.Complex128,
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
