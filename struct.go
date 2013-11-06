package fat

import (
	"fmt"
	"github.com/metakeule/meta"
	"reflect"
	"strings"
)

func StructType(østruct interface{}) string {
	ty := reflect.TypeOf(østruct).Elem()
	return "*" + ty.PkgPath() + "." + ty.Name()
}

// sets all attributes of a struct that are of type *Field
// to a *Field with the Type set to what is given in the tag "fat.type"
// with defaults set to what is given in the tag "fat.default"
// and with enums set to what is in the tag "fat.enum", separated by pipe symbols (|)
func Proto(østruct interface{}) (østru interface{}) {
	structtype := StructType(østruct)
	fn := func(field reflect.StructField, val reflect.Value) {
		if _, ok := val.Interface().(*Field); ok {
			ty := field.Tag.Get("fat.type")
			if ty == "" {
				panic(fmt.Sprintf("struct %s has no fat.type tag for field %s", structtype, field.Name))
			}

			f := newSpec(structtype, field.Name, ty)
			def := field.Tag.Get("fat.default")
			if def != "" {
				d := f.fieldSpec.new()
				err := d.Scan(def)
				if err != nil {
					panic(fmt.Sprintf("default value %s for field %s in struct %s is not of type %s",
						def, field.Name, structtype, d.Typ()))
				}
				f.default_ = d
			}
			enum := field.Tag.Get("fat.enum")
			if enum != "" {
				enumVals := strings.Split(enum, "|")
				enums := make([]Type, len(enumVals))
				for i, en := range enumVals {
					e := f.fieldSpec.new()
					err := e.Scan(en)
					if err != nil {
						panic(fmt.Sprintf("enum value %s for field %s in struct %s is not of type %s",
							en, field.Name, structtype, e.Typ()))
					}
					enums[i] = e
				}
				f.enum = enums
			}
			val.Set(reflect.ValueOf(f))
		}
	}
	meta.Struct.EachRaw(østruct, fn)
	return østruct
}

// prefills the given newstruct based on the given prototype
func New(øprototype interface{}, ønewstruct interface{}) (ønew interface{}) {
	prototype := reflect.TypeOf(øprototype).String()
	newtype := reflect.TypeOf(ønewstruct).String()
	if prototype != newtype {
		panic(fmt.Sprintf("prototype (%s) does and new type (%s) are not the same", prototype, newtype))
	}
	proto := reflect.ValueOf(øprototype).Elem()
	fn := func(field reflect.StructField, val reflect.Value) {
		if _, ok := val.Interface().(*Field); ok {
			val.Set(
				reflect.ValueOf(
					proto.FieldByName(field.Name).
						Interface().(*Field).
						New(),
				),
			)
		}
	}
	meta.Struct.EachRaw(ønewstruct, fn)
	return ønewstruct
}
