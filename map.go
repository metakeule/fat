package fat

import (
	"encoding/json"
	"fmt"
	"time"
)

type map_ struct {
	typ string
	Map map[string]Type
}

func (ømap *map_) Typ() string      { return "map." + ømap.typ }
func (ømap *map_) Get() interface{} { return ømap.Map }

func (ømap *map_) Set(i interface{}) error {
	_map, ok := i.(*map_)
	if !ok {
		return fmt.Errorf("can't convert %T to map.%s", i, ømap.typ)
	}
	if ømap.typ != _map.typ {
		return fmt.Errorf("can't set map.%s to map of type map.%s", ømap.typ, _map.typ)
	}
	ømap.Map = _map.Map
	return nil
}

func (ømap *map_) String() string {
	var data []byte
	var err error
	if ømap.typ == "time" {
		asStrings := make(map[string]string, len(ømap.Map))
		for k, t := range ømap.Map {
			asStrings[k] = t.String()
		}
		data, err = json.Marshal(asStrings)
	} else {
		data, err = json.Marshal(ømap.Map)
	}
	if err != nil {
		panic(fmt.Sprintf("can't convert map.%s to string: %s", ømap.typ, err.Error()))
	}
	return string(data)
}
func (ømap *map_) Scan(n string) error {
	_map := map[string]interface{}{}
	err := json.Unmarshal([]byte(n), &_map)
	if err != nil {
		return err
	}
	ømap.Map = map[string]Type{}
	for k, v := range _map {
		t := newType(ømap.typ)
		var e error
		switch ømap.typ {
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
		ømap.Map[k] = t
	}
	return nil
}

func mustSetTypeForMap(_map map[string]Type, typ string, k string, v interface{}) {
	mustBeUTF8(k)
	t := newType(typ)
	e := t.Set(v)
	if e != nil {
		panic(fmt.Sprintf("can't set %#v (%T) to %s", v, v, typ))
	}
	_map[k] = t
}

func Map(givenmap interface{}) *map_ {
	_map := map[string]Type{}
	var typ string
	switch m := givenmap.(type) {
	case map[string]string:
		typ = "string"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]int8:
		typ = "int"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]int16:
		typ = "int"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]int32:
		typ = "int"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]int64:
		typ = "int"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]int:
		typ = "int"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]float64:
		typ = "float"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]float32:
		typ = "float"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]bool:
		typ = "bool"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	case map[string]time.Time:
		typ = "time"
		for k, v := range m {
			mustSetTypeForMap(_map, typ, k, v)
		}
	default:
		panic(fmt.Sprintf("unsupported type: %T", givenmap))
	}
	return &map_{typ: typ, Map: _map}
}

var mapTypes = map[string]bool{
	"string": true,
	"int":    true,
	"float":  true,
	"bool":   true,
	"time":   true,
}

// panics if there are different types
func MapType(typ string, vals ...interface{}) *map_ {
	mustBeUTF8(typ)
	if !mapTypes[typ] {
		panic("unsupported type: " + typ)
	}
	if len(vals)%2 != 0 {
		panic(fmt.Sprintf("map must be given pairs of string interface, len is odd: %v", vals))
	}
	_map := map[string]Type{}
	for i := 0; i < len(vals)-1; i += 2 {
		k, ok := vals[i].(string)
		if !ok {
			panic(fmt.Sprintf("is no string %#v (%T)", vals[i], vals[i]))
		}
		mustBeUTF8(k)
		v := vals[i+1]
		t := newType(typ)
		e := t.Set(v)
		if e != nil {
			panic(fmt.Sprintf("can't set %#v (%T) to %s", v, v, typ))
		}
		_map[k] = t
	}
	return &map_{typ: typ, Map: _map}
}

func MapStrings(m map[string]string) *map_  { return Map(m) }
func MapInts(m map[string]int64) *map_      { return Map(m) }
func MapFloats(m map[string]float64) *map_  { return Map(m) }
func MapBools(m map[string]bool) *map_      { return Map(m) }
func MapTimes(m map[string]time.Time) *map_ { return Map(m) }
