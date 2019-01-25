# bitfield (WIP)
struct bit fields + un(packing), similar-ish to C bit fields (with Go idioms)

Example:
```
type s1 struct {
    X bool      `len:"1"`
    Y uint      `len:"8"` // bits length
    _ struct{}  `len:"0"` // fill to nearest byte
    Z []s2
}

type s2 struct {
    X []byte    `len:"8"` // bytes length
    Y []byte    `upack:"customUnpack" pack:"customPack"`
    Z []byte    `len:"0"` // fill all the bytes left
}

dst := s1{}
bitfield.Unpack(&dst, []byte{})
```

Notice:
```
bitfield.Unpack(&dst, src) is equivalence of bitfield.BigEndian.Unpack(&dst, src)
```
For little endianness, use:
```
bitfield.LittleEndian.Unpack(&dst, src)
```