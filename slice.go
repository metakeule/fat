package fat

import (
	"encoding/json"
	"fmt"
	"time"
)

// for more complex types
type slice struct {
	typ   string
	Slice []Type
}

func (øslice *slice) Typ() string      { return "slice." + øslice.typ }
func (øslice *slice) Get() interface{} { return øslice.Slice }

func (øslice *slice) Set(i interface{}) error {
	s, ok := i.(*slice)
	if !ok {
		return fmt.Errorf("can't convert %T to slice.%s", i, øslice.typ)
	}
	if øslice.typ != s.typ {
		return fmt.Errorf("can't set slice.%s to slice of type slice.%s", øslice.typ, s.typ)
	}
	øslice.Slice = s.Slice
	return nil
}

func (øslice *slice) String() string {
	var data []byte
	var err error
	if øslice.typ == "time" {
		asStrings := make([]string, len(øslice.Slice))
		for i, t := range øslice.Slice {
			asStrings[i] = t.String()
		}
		data, err = json.Marshal(asStrings)
	} else {
		data, err = json.Marshal(øslice.Slice)
	}
	if err != nil {
		panic(fmt.Sprintf("can't convert slice.%s to string: %s", øslice.typ, err.Error()))
	}
	return string(data)
}

func (øslice *slice) Scan(s string) error {
	intfSlice := []interface{}{}
	err := json.Unmarshal([]byte(s), &intfSlice)
	if err != nil {
		return err
	}
	øslice.Slice = make([]Type, len(intfSlice))

	for i, v := range intfSlice {
		t := newType(øslice.typ)
		var e error
		switch øslice.typ {
		case "int":
			switch vt := v.(type) {
			case float64:
				e = t.Set(int64(vt))
			case float32:
				e = t.Set(int64(vt))
			default:
				e = fmt.Errorf("can't convert %#v (%T) to int", v)
			}
		case "time":
			e = t.Scan(v.(string))
		default:
			e = t.Set(v)
		}

		if e != nil {
			return e
		}
		øslice.Slice[i] = t
	}
	return nil
}

// panics if there are different types
func Slice(typ string, intfs ...interface{}) *slice {
	types := make([]Type, len(intfs))
	for i, intf := range intfs {
		t := newType(typ)
		err := t.Set(intf)
		if err != nil {
			panic(fmt.Sprintf("can't set value to type %s at position %v: %s", typ, i, err.Error()))
		}
		types[i] = t
	}
	return &slice{typ: typ, Slice: types}
}

func Strings(strings ...string) *slice {
	params := make([]interface{}, len(strings))
	for i, s := range strings {
		params[i] = s
	}
	return Slice("string", params...)
}
func Ints(params ...interface{}) *slice {
	return Slice("int", params...)
}
func Floats(params ...interface{}) *slice {
	return Slice("float", params...)
}
func Bools(bools ...bool) *slice {
	params := make([]interface{}, len(bools))
	for i, b := range bools {
		params[i] = b
	}
	return Slice("bool", params...)
}
func Times(times ...time.Time) *slice {
	params := make([]interface{}, len(times))
	for i, t := range times {
		params[i] = t
	}
	return Slice("time", params...)
}
