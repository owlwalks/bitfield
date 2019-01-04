# bitfield (WIP)
struct bit fields + un(packing), similar-ish to C bit fields (with Go idioms)

Example:
```
type A struct {
    x bool      `len:"1"`
    y uint      `len:"8"` // bits length
    _ struct{}  `len:"0"` // fill to nearest byte
    z []B
}

type B struct {
    x []byte    `len:"8"` // bytes length
    y []byte    `upack:"customUnpack" pack:"customPack"`
    z []byte    `len:"0"` // fill all the bytes left
}
```
