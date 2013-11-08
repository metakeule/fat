package fat

import (
	"fmt"
	. "github.com/metakeule/nil"
)

type fieldSpec struct {
	// constructs a new value
	new              func() Type
	name, structtype string
	// default value
	default_ Type
	// list of valid entries
	enum []Type
	// Validator function, returns all errors
	Validator func(øself *Field) []error
}

func (øfieldSpec *fieldSpec) Name() string       { return øfieldSpec.name }
func (øfieldSpec *fieldSpec) Default() Type      { return øfieldSpec.default_ }
func (øfieldSpec *fieldSpec) Enum() []Type       { return øfieldSpec.enum }
func (øfieldSpec *fieldSpec) StructType() string { return øfieldSpec.structtype }
func (øfieldSpec *fieldSpec) Path() string {
	return øfieldSpec.structtype + "." + øfieldSpec.name
}

func (øfieldSpec *fieldSpec) New(østruct interface{}) *Field {
	return &Field{
		Type:            øfieldSpec.new(),
		Struct:          østruct,
		fieldSpec:       øfieldSpec,
		IsSet:           false,
		FailedScanInput: "",
	}
}

// creates a new Field based on a spec that is created by the way
// the first of the given vals is considered as default value (nil is considered
// to have no default), the rest are enums
func newSpec(structtype string, østruct interface{}, field string, typ string, vals ...Type) *Field {
	var default_ Type
	var enum []Type
	if len(vals) > 0 {
		default_ = vals[0]
		enum = append(enum, vals[1:]...)
	}
	newFunc, ok := newTypeFuncs[typ]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", typ))
	}
	return genSpec(
		structtype,
		field,
		newFunc,
		default_,
		enum,
	).New(østruct)
}

func genSpec(structtype string, field string, øtypeGenerator func() Type, default_ Type, enum []Type) (fs *fieldSpec) {
	Nø(øtypeGenerator)
	fs = &fieldSpec{
		name:       field,
		structtype: structtype,
		new:        øtypeGenerator,
		default_:   default_,
		enum:       enum,
	}
	return
}
