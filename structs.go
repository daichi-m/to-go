package togo

import (
	"fmt"
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

// Primitive checks if an instance of FieldDT is a primitive type
func (f FieldDT) Primitive() bool {
	switch f {
	case Int, Bool, Int64, Float64, String:
		return true
	default:
		return false
	}
}

func (f FieldDT) String() string {
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

func (f FieldDT) GoString() string {
	return fmt.Sprintf("FieldDT: %s", f.String())
}

// NewFieldDT converts any instance of interface{} type to FieldDT
func NewFieldDT(val interface{}) (FieldDT, bool) {
	k := reflect.ValueOf(val).Kind()
	return fieldDT(k)
}

// fieldDT is the private method to convert a reflect.Kind into FieldDT
func fieldDT(k reflect.Kind) (FieldDT, bool) {
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

/*
// Cloner interface defines a clone method that can create a cloned copy of itself
type Cloner interface {
	Clone() (Cloner, error)
}

// Grower interface grows the object to accomodate further changes from another Grower
type Grower interface {
	Grow(other Grower) (Grower, error)
}

type Equaler interface {
	Equals(other Equaler) bool
}*/

// IField is a cover interface for a Field in a GoStruct.
type IField interface {
	fmt.Stringer
	fmt.GoStringer

	Clone() IField
	Grow(other IField) (IField, error)
	Equals(other IField) bool
	Name() string

	/*
		// Getters for IField

		Annotation() string
		DataType() FieldDT
		DTStruct() string
		SliceNesting() int

		// Setters for IField
		SetName(name string) IField
		SetAnnotation(annotation string) IField
		SetDataType(dt FieldDT) IField

	*/
}

// IFieldMaker is a factory interface for IField
type IFieldMaker interface {
	MakeIField(name string, annotation string, dataType FieldDT,
		dts string, nest int) IField
}

// FieldMaker is a IFieldMaker for the *Field type
type FieldMaker struct{}

var _ IFieldMaker = FieldMaker{}

// Field is the struct that donates a go Field inside the go struct that will
// be generated. It has a Name, a Type and optionally an
// Annotation to user for (un)marshalling
type Field struct {
	FieldName    string
	Annotation   string
	DataType     FieldDT
	DTStruct     string
	SliceNesting int
}

var _ IField = (*Field)(nil)

// Name returns the name of the field
func (f *Field) Name() string {
	return f.FieldName
}

func (f *Field) String() string {

	nest := strings.Builder{}
	if f.SliceNesting > 0 {
		nest.WriteString("[]")
	}

	if f.DataType.Primitive() {
		return fmt.Sprintf("%s\t %s %s\t `%s`",
			f.FieldName, nest.String(), f.DataType.String(), f.Annotation)
	}
	return fmt.Sprintf("%s\t %s %s\t `%s`",
		f.FieldName, nest.String(), f.DTStruct, f.Annotation)
}

func (f *Field) GoString() string {
	return fmt.Sprintf("Field: %s", f.String())
}

// Equals check if this instance of field is "in-principle"
// equal to other instance
func (f *Field) Equals(eq IField) bool {

	of, ok := eq.(*Field)
	if !ok {
		return false
	}

	if f.FieldName != of.FieldName {
		return false
	}
	if f.DataType != of.DataType {
		return false
	}
	if f.DTStruct != of.DTStruct {
		return false
	}
	if f.SliceNesting != of.SliceNesting {
		return false
	}
	return true
}

// Grow adds annotation from the second IField into this Field
func (f *Field) Grow(oif IField) (IField, error) {
	other, ok := oif.(*Field)
	if !ok {
		return nil, fmt.Errorf("Cannot grow from an non IField type")
	}

	if len(f.Annotation) == 0 {
		f.Annotation = other.Annotation
	}

	if !strings.Contains(f.Annotation, other.Annotation) {
		f.Annotation = fmt.Sprintf("%s,%s", f.Annotation, other.Annotation)
	}
	return f, nil
}

// Clone a Field into another Field object
func (f *Field) Clone() IField {
	of := Field{
		FieldName:    f.FieldName,
		Annotation:   f.Annotation,
		DataType:     f.DataType,
		DTStruct:     f.DTStruct,
		SliceNesting: f.SliceNesting,
	}
	return &of
}

// MakeIField converts the generic interface{} into a IField with the given name.
// It throws an error in case the field creation is not successful due to some reason.
func (fm FieldMaker) MakeIField(name string, annotation string, dataType FieldDT,
	dts string, nest int) IField {

	f := Field{
		FieldName:    name,
		Annotation:   annotation,
		DataType:     dataType,
		DTStruct:     dts,
		SliceNesting: nest,
	}
	return &f
}

// IGoStruct is a cover interface for a GoStruct
type IGoStruct interface {
	fmt.Stringer
	fmt.GoStringer

	Clone() IGoStruct
	Grow(other IGoStruct) (IGoStruct, error)
	Equals(other IGoStruct) bool
	AddField(field IField) (IGoStruct, error)
	Name() string
	Level() int
	ChangeName(string)

	// Name() string
	// Fields() []IField
	// Level() int

	// SetName(name string) IGoStruct
	// SetFields(fields []IField) IGoStruct
	// SetLevel(level int) IGoStruct

}

// GoStruct is the representation of a go struct. It has a name,
// a set of Field types and a Level to determine at what level should the
// struct be defined in the final source code.
type GoStruct struct {
	StructName  string
	Fields      map[string]IField
	StructLevel int
}

var _ IGoStruct = (*GoStruct)(nil)

// Name returns the name of the GoStruct
func (gs *GoStruct) Name() string {
	return gs.StructName
}

// Level returns the level of the GoStruct
func (gs *GoStruct) Level() int {
	return gs.StructLevel
}

// ChangeName changes the name of the GoStruct
func (gs *GoStruct) ChangeName(name string) {
	gs.StructName = name
}

func (gs *GoStruct) String() string {
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("type %s struct", gs.StructName))
	for _, f := range gs.Fields {
		sb.WriteString(fmt.Sprintf("\n \t %s", f.String()))
	}
	sb.WriteString(fmt.Sprintf("\n At Level %d", gs.StructLevel))
	return sb.String()
}

func (gs *GoStruct) GoString() string {
	return gs.String()
}

// Clone deep clones a GoStruct. Visible for testing
func (gs *GoStruct) Clone() IGoStruct {
	ngs := GoStruct{
		StructName:  gs.StructName,
		Fields:      make(map[string]IField),
		StructLevel: gs.StructLevel,
	}
	for n, f := range gs.Fields {
		ngs.Fields[n] = f.Clone()
	}
	return &ngs
}

// AddField adds a field to the GoStruct instance
func (gs *GoStruct) AddField(ifl IField) (IGoStruct, error) {
	f := ifl.(*Field)
	logger := getLogger().Sugar()
	defer logger.Sync()

	if reflect.ValueOf(f).IsNil() {
		return nil, fmt.Errorf("Cannot add a nil field")
	}

	if gs.Fields == nil {
		gs.Fields = make(map[string]IField)
	}
	exFld, ok := gs.Fields[f.FieldName]
	if !ok {
		gs.Fields[f.FieldName] = f
		return gs, nil
	}

	if !exFld.Equals(f) {
		return nil, fmt.Errorf("Unmatched field, expected %v, found %v", exFld, f)
	}
	f.Grow(exFld)
	gs.Fields[f.FieldName] = f
	logger.Debugf("Added field %s to the GoStruct %s", f.FieldName, gs.StructName)
	return gs, nil
}

// Equals check if two GoStruct instances are equal.
// Equality of GoStructs depends solely on name.
// Fields can get added and deleted, so field equality is not checked
func (gs *GoStruct) Equals(igs IGoStruct) bool {

	if reflect.ValueOf(igs).IsNil() {
		return false
	}
	other, ok := igs.(*GoStruct)
	if !ok {
		return false
	}
	if gs.StructName == other.StructName {
		return true
	}
	return false
}

// Grow a struct with additional fields from the other GoStruct instance
func (gs *GoStruct) Grow(og IGoStruct) (IGoStruct, error) {

	logger := getLogger().Sugar()
	defer logger.Sync()

	var other *GoStruct
	var ok bool

	if other, ok = og.(*GoStruct); !ok {
		return nil, fmt.Errorf("GoStruct: Cannot grow from different object")
	}

	if eq := gs.Equals(other); !eq {
		logger.Debugf("Structs %s and %s are not equal, cannot grow", gs.StructName, other.StructName)
		return nil, fmt.Errorf("GoStruct %s is not equal, cannot grow", other.StructName)
	}

	for name, field := range other.Fields {
		gfield, ok := gs.Fields[name]
		if !ok {
			gs.AddField(field)
			continue
		}
		if feq := gfield.Equals(field); !feq {
			return nil, fmt.Errorf("Field %s and %s are not equal, cannot grow",
				field, gfield)
		}
	}
	return gs, nil
}

// IsEmpty checks if this GoStruct is empty
// (i.e., this instance does not have a name)
/*
func (gs *GoStruct) isEmpty() bool {
	if len(gs.name) == 0 {
		return true
	}
	return false
}*/

// IGoStructMaker is a factory interface for IGoStruct
type IGoStructMaker interface {
	MakeGoStruct(name string, fields []IField, level int) IGoStruct
}

// GoStructMaker is a IGoStructMaker for the *GoStruct type
type GoStructMaker struct{}

var _ IGoStructMaker = GoStructMaker{}

// MakeGoStruct makes an instance of GoStruct from the corresponding name, fields and level values
func (gsm GoStructMaker) MakeGoStruct(name string, fields []IField, level int) IGoStruct {
	gs := GoStruct{
		StructName:  name,
		StructLevel: level,
		Fields:      make(map[string]IField),
	}
	if fields != nil {
		for _, f := range fields {
			f2 := f.(*Field)
			gs.Fields[f2.FieldName] = f2
		}
	}
	return &gs
}
