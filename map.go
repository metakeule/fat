package fat

import (
	"encoding/json"
	"fmt"
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

// panics if there are different types
func Map(typ string, vals ...interface{}) *map_ {
	mustBeUTF8(typ)
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

func MapStrings(params ...interface{}) *map_ { return Map("string", params...) }
func MapInts(params ...interface{}) *map_    { return Map("int", params...) }
func MapFloats(params ...interface{}) *map_  { return Map("float", params...) }
func MapBools(params ...interface{}) *map_   { return Map("bool", params...) }
func MapTimes(params ...interface{}) *map_   { return Map("time", params...) }
