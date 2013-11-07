package fat

import (
	"encoding/json"
	"fmt"
)

type Field struct {
	Type
	*fieldSpec
	IsSet bool // is true, if the value was set, i.e. the type was correct, it may however by invalid
	// saves the input for a failed scan,
	// is empty if the scan did not fail
	FailedScanInput string
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
// field may be set and invalid at the same time
// IsSet only tells us, wether the type is correct
func (øField *Field) Set(i interface{}) error {
	err := øField.Type.Set(i)
	if err != nil {
		return err
	}
	øField.IsSet = true
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

func (øField *Field) MarshalJSON() ([]byte, error) {
	switch øField.Type.(type) {
	case *slice:
		return []byte(øField.String()), nil
	case *map_:
		return []byte(øField.String()), nil
	case *time_:
		return json.Marshal(øField.String())
	default:
		return json.Marshal(øField.Get())
	}
}

func (øField *Field) UnmarshalJSON(data []byte) (err error) {
	var target string
	switch øField.Type.(type) {
	case *string_, *time_:
		err = json.Unmarshal(data, &target)
	default:
		target = string(data)
	}
	if err == nil {
		err = øField.Scan(target)
	}
	return
}

// overwrite Type.String to return default value, if IsSet is false
func (øField *Field) String() string {
	if øField.default_ != nil && !øField.IsSet {
		return øField.default_.String()
	}
	return øField.Type.String()
}

// overwrite Type.Scan to track, if the field was set
// field only is set if scan was successful
// Scan does not validation check, that must be run after
// Scan, or use ScanAndValidate
func (øField *Field) Scan(s string) error {
	if err := shouldBeUTF8(s); err != nil {
		return err
	}
	øField.FailedScanInput = s
	err := øField.Type.Scan(s)
	if err != nil {
		return err
	}
	øField.IsSet = true
	øField.FailedScanInput = ""
	return nil
}

func (øField *Field) ScanAndValidate(s string) (errs []error) {
	errs = []error{}
	err := øField.Scan(s)
	if err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, øField.Validate()...)
	return
}

// calls Scan and panics if there is an error
func (øField *Field) MustScan(s string) {
	err := øField.Scan(s)
	if err != nil {
		panic(err.Error())
	}
}

func (øField *Field) MustScanAndValidate(s string) {
	øField.MustScan(s)
	øField.MustValidate()
}

// overwrite Type.Get to return default value, if IsSet is false
func (øField *Field) Get() (i interface{}) {
	if øField.default_ != nil && !øField.IsSet {
		return øField.default_.Get()
	}
	return øField.Type.Get()
}

// panics on first error
func (øField *Field) MustValidate() {
	errs := øField.Validate()
	if len(errs) > 0 {
		panic(errs[0])
	}
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
