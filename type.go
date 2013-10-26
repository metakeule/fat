package fat

import (
	"fmt"
)

type Type interface {
	Typ() string
	Get() interface{}  // must be typecasted
	String() string    // represent value as string
	Scan(string) error // sets the value based on a string
	Set(interface{}) error
}

// maps type to a generator
var newTypeFuncs = map[string]func() Type{
	"int":    func() Type { return Int(0) },
	"float":  func() Type { return Float(float64(0.0)) },
	"bool":   func() Type { return Bool(false) },
	"string": func() Type { return String("") },
	"time":   func() Type { return Time(zeroTime) },

	"map.int":    func() Type { return newMap("int") },
	"map.float":  func() Type { return newMap("float") },
	"map.bool":   func() Type { return newMap("bool") },
	"map.string": func() Type { return newMap("string") },
	"map.time":   func() Type { return newMap("time") },

	"slice.int":    func() Type { return newSlice("int") },
	"slice.float":  func() Type { return newSlice("float") },
	"slice.bool":   func() Type { return newSlice("bool") },
	"slice.string": func() Type { return newSlice("string") },
	"slice.time":   func() Type { return newSlice("time") },
}

func newType(typ string) Type {
	f, ok := newTypeFuncs[typ]
	if !ok {
		panic(fmt.Sprintf("can't create value of type: %s: unknown type", typ))
	}
	return f()
}

func newMap(typ string) Type   { return &map_{typ: typ, Map: map[string]Type{}} }
func newSlice(typ string) Type { return &slice{typ: typ, Slice: []Type{}} }
