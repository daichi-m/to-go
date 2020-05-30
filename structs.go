package togo

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// FieldDT is an alias for the field data type represented
// as an uint for convenience purpose
type FieldDT uint

// Constant describing the different the field data type supported currently.
// Some assumptions made here are:
//   1. Smaller int's (int8, int16) is widened to int32 in the struct.
//   2. Similar float32 is widened to float64 in the struct.
//   3. Two complex data type supported are Slice and Map. All maps and slices would be
//		represented by these two consts always.
const (
	Initial = iota
	Bool
	Int
	Int64
	Float64
	String
	Slice
	Map
)

func (f FieldDT) primitive() bool {
	switch f {
	case Int, Bool, Int64, Float64, String:
		return true
	default:
		return false
	}
}

func (f FieldDT) str() string {
	switch f {
	case Initial:
		return "Initial"
	case Bool:
		return "Bool"
	case Int:
		return "Int"
	case Int64:
		return "Int64"
	case Float64:
		return "Float64"
	case String:
		return "String"
	case Slice:
		return "Slice"
	case Map:
		return "Map"
	default:
		return "Unknown Type"
	}
}

func toFieldDT(k reflect.Kind) (FieldDT, bool) {
	if k == reflect.Int {
		return Int, true
	} else if k == reflect.Int64 {
		return Int64, true
	} else if k == reflect.Float32 || k == reflect.Float64 {
		return Float64, true
	} else if k == reflect.String {
		return String, true
	} else if k == reflect.Map {
		return Map, true
	} else if k == reflect.Slice {
		return Slice, true
	} else if k == reflect.Bool {
		return Bool, true
	} else {
		return Initial, false
	}
}

// Field is the struct that donates a go Field inside the go struct that will
// be generated. It has a Name, a Type and optionally an
// Annotation to user for (un)marshalling
type Field struct {
	name         string
	annotation   string
	dataType     FieldDT
	dtStruct     string
	sliceNesting int
}

// Equals check if this instance of field is "in-principle"
// equal to other instance
func (f *Field) Equals(of *Field) bool {
	if f.name != of.name {
		return false
	}
	if f.dataType != of.dataType {
		return false
	}
	if f.dtStruct != of.dtStruct {
		return false
	}
	if f.sliceNesting != of.sliceNesting {
		return false
	}
	return true
}

// Annotate adds an annotation to the field
func (f *Field) Annotate(a string) {
	if f.annotation == "" {
		f.annotation = a
	}

	if !strings.Contains(f.annotation, a) {
		f.annotation = fmt.Sprintf("%s,%s", f.annotation, a)
	}
}

// Clones a field. Visible for testing
func (f *Field) clone() Field {
	return Field{
		f.name, f.annotation, f.dataType, f.dtStruct, f.sliceNesting,
	}
}

// ToField converts the generic interface{} into a Field with the given name.
// It throws an error in case the field creation is not successful due to some reason.
func ToField(name string, val interface{}) (*Field, error) {
	f := new(Field)
	// TODO: To work on normalizing the name
	f.name = name
	k := reflect.ValueOf(val).Kind()
	dt, ok := toFieldDT(k)
	if !ok {
		return nil, UnsupportedType{
			data: val,
		}
	}
	f.dataType = dt
	// TODO: Revisit this - need to fill up the annotation, slice nesting and dtstruct
	f.annotation = ""
	log.Printf("Field created: %+v \n", f)
	return f, nil
}

// GoStruct is the representation of a go struct. It has a name,
// a set of Field types and a Level to determine at what level should the
// struct be defined in the final source code.
type GoStruct struct {
	Name   string
	Fields map[string]*Field
	Level  int
}

// Clone deep clones a GoStruct. Visible for testing
func (gs GoStruct) clone() GoStruct {
	ngs := GoStruct{
		Name:   gs.Name,
		Fields: make(map[string]*Field),
		Level:  gs.Level,
	}
	for n, f := range gs.Fields {
		nf := f.clone()
		ngs.Fields[n] = &nf
	}
	return ngs
}

// AddField adds a field to the GoStruct instance
func (gs *GoStruct) AddField(f *Field) error {
	if f == nil {
		return GoStructError{
			gs:      *gs,
			message: "Attempt to add a nil Field",
		}
	}

	if gs.Fields == nil {
		gs.Fields = make(map[string]*Field)
	}
	exFld, ok := gs.Fields[f.name]
	if !ok {
		gs.Fields[f.name] = f
		return nil
	}

	if !exFld.Equals(f) {
		return GoStructError{
			gs:      *gs,
			message: fmt.Sprintf("Unmatched Fields. Have %+v, received: %+v", exFld, f),
		}
	}
	f.Annotate(exFld.annotation)
	gs.Fields[f.name] = f
	log.Printf("Added field %+v to the GoStruct %+v", f.name, gs.Name)
	return nil
}

// Equals check if two GoStruct instances are equal.
// Equality of GoStructs depends solely on name.
// Fields can get added and deleted, so field equality is not checked
func (gs *GoStruct) Equals(other *GoStruct) bool {
	if other == nil {
		return false
	}
	if gs.Name == other.Name {
		return true
	}
	return false
}

// Grow a struct with additional fields from the other GoStruct instance
func (gs *GoStruct) Grow(other *GoStruct) error {
	if eq := gs.Equals(other); !eq {
		log.Printf("Structs %+v and %+v are not equal, cannot grow", gs, other)
		return GoStructError{
			gs:      *gs,
			message: fmt.Sprintf("Struct %s not equal, cannot Grow", other.Name),
		}
	}

	for name, field := range other.Fields {
		gfield, ok := gs.Fields[name]
		if !ok {
			gs.AddField(field)
			continue
		}
		if feq := gfield.Equals(field); !feq {
			return GoStructError{
				gs:      *gs,
				message: fmt.Sprintf("Field %s does not equal, cannot Grow", gfield.name),
			}
		}
	}
	return nil
}

// IsEmpty checks if this GoStruct is empty
// (i.e., this instance does not have a name)
func (gs *GoStruct) IsEmpty() bool {
	if len(gs.Name) == 0 {
		return true
	}
	return false
}
