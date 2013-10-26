package fat

import (
	"fmt"
	"strconv"
)

type bool_ bool

func Bool(b bool) *bool_         { b_ := bool_(b); return &b_ }
func (b bool_) Typ() string      { return "bool" }
func (b bool_) Get() interface{} { return bool(b) }
func (b bool_) String() string   { return fmt.Sprintf("%v", bool(b)) }

func (øbool *bool_) Set(i interface{}) error {
	b, ok := i.(bool)
	if !ok {
		return fmt.Errorf("can't convert %T to bool", i)
	}
	*øbool = bool_(b)
	return nil
}

func (øbool *bool_) Scan(s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*øbool = bool_(b)
	return nil
}
