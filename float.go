package fat

import (
	"fmt"
	"strconv"
)

type float_ float64

func Float(f float64) *float_     { f_ := float_(f); return &f_ }
func (f float_) Typ() string      { return "float" }
func (f float_) Get() interface{} { return f }
func (f float_) String() string   { return fmt.Sprintf("%v", float64(f)) }

func (øfloat *float_) Scan(s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*øfloat = float_(f)
	return nil
}

func (øfloat *float_) Set(i interface{}) error {
	var f float64
	switch t := i.(type) {
	case float32:
		f = float64(t)
	case float64:
		f = t
	default:
		return fmt.Errorf("can't convert %T to float", i)
	}
	*øfloat = float_(f)
	return nil
}
