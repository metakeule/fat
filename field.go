package fat

import (
	"fmt"
)

type Field struct {
	Type
	*fieldSpec
	IsSet bool // is true, if the value was set, may be faked
	// not sure if this is neccessary
	// Struct                interface{}
}

// calls Set and panics if there is an error
func (øField *Field) MustSet(i interface{}) {
	err := øField.Set(i)
	if err != nil {
		panic(err.Error())
	}
}

// overwrite Type.Set to track, if the field was set
func (øField *Field) Set(i interface{}) error {
	øField.IsSet = true

	err := øField.Type.Set(i)
	if err != nil {
		return err
	}
	if !øField.IsValid() {
		return fmt.Errorf("validation errors, run Validate() for specific errors")
	}
	return nil
}

// returns if  field is valid
func (øField *Field) IsValid() bool {
	errs := øField.Validate()
	return !(len(errs) > 0)
}

// overwrite Type.String to return default value, if IsSet is false
func (øField *Field) String() string {
	if øField.default_ != nil && !øField.IsSet {
		return øField.default_.String()
	}
	return øField.Type.String()
}

// overwrite Type.Scan to track, if the field was set
func (øField *Field) Scan(s string) error {
	øField.IsSet = true
	err := øField.Type.Scan(s)
	if err != nil {
		return err
	}
	if !øField.IsValid() {
		return fmt.Errorf("validation errors, run Validate() for specific errors")
	}
	return nil
}

// calls Scan and panics if there is an error
func (øField *Field) MustScan(s string) {
	err := øField.Scan(s)
	if err != nil {
		panic(err.Error())
	}
}

// overwrite Type.Get to return default value, if IsSet is false
func (øField *Field) Get() (i interface{}) {
	if øField.default_ != nil && !øField.IsSet {
		return øField.default_.Get()
	}
	return øField.Type.Get()
}

// validates the content of field and returns all errors
func (øField *Field) Validate() (errs []error) {
	if len(øField.enum) > 0 {
		var valid bool
		var val = øField.Get()
		for _, en := range øField.enum {
			if val == en.Get() {
				valid = true
				break
			}
		}
		if !valid {
			errs = append(errs, fmt.Errorf("value not part of enum values"))
		}
	}
	if øField.Validator != nil {
		errs = append(errs, øField.Validator(øField)...)
	}
	return
}
