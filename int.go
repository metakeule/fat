package fat

import (
	"fmt"
	"strconv"
)

type int_ int64

func Int(i int64) *int_         { i_ := int_(i); return &i_ }
func (i int_) Typ() string      { return "int" }
func (i int_) Get() interface{} { return int64(i) }
func (i int_) String() string   { return fmt.Sprintf("%v", int64(i)) }

func (øi *int_) Scan(s string) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*øi = int_(i)
	return nil
}

func (øi *int_) Set(i interface{}) error {
	var i_ int64
	switch t := i.(type) {
	case int8:
		i_ = int64(t)
	case int16:
		i_ = int64(t)
	case int32:
		i_ = int64(t)
	case int64:
		i_ = t
	case int:
		i_ = int64(t)
	default:
		return fmt.Errorf("can't convert %T to int", i)
	}

	*øi = int_(i_)
	return nil
}
