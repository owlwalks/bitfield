# bitfield (WIP)
WIP: struct bit fields + un(packing), similar-ish to C bit fields (with Go idioms)

Example:
```
type A struct {
    x bool      `len:"1"`
    y uint      `len:"4"`
    _ struct{}  `len:"0"` // fill to nearest byte
    z []B
}

type B struct {
    x int       `len:"8"`
    y []byte    `upack:"customUnpack" pack:"customPack"`
}
```
