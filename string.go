package fat

import (
	"fmt"
	"unicode/utf8"
)

func shouldBeUTF8(s string) error {
	if !utf8.ValidString(s) {
		return fmt.Errorf("%+q is not a valid utf8 string", s)
	}
	return nil
}

func mustBeUTF8(s string) {
	err := shouldBeUTF8(s)
	if err != nil {
		panic(err.Error())
	}
}

// represents a string that is valid utf8
type string_ string

func String(s string) *string_ {
	mustBeUTF8(s)
	s_ := string_(s)
	return &s_
}

func (s string_) Typ() string      { return "string" }
func (s string_) Get() interface{} { return string(s) }
func (s string_) String() string   { return string(s) }
func (østring *string_) Scan(s string) error {
	*østring = string_(s)
	return nil
}

func (østring *string_) Set(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("can't convert %T to string", i)
	}
	if err := shouldBeUTF8(s); err != nil {
		return err
	}
	*østring = string_(s)
	return nil
}
