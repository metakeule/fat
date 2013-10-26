package fat

import (
	"fmt"
)

type string_ string

func String(s string) *string_ { s_ := string_(s); return &s_ }

func (s string_) Typ() string                { return "string" }
func (s string_) Get() interface{}           { return string(s) }
func (s string_) String() string             { return string(s) }
func (østring *string_) Scan(s string) error { *østring = string_(s); return nil }

func (østring *string_) Set(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("can't convert %T to string", i)
	}
	*østring = string_(s)
	return nil
}
