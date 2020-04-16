package togo

import (
	"errors"
	"fmt"
	"strings"
)

// The struct that donates a go field inside the go struct that will be generated.
// It has a Name, a Type and optionally an Annotation to user for (un)marshalling
type Field struct {
	Name, Annotation, TypeString string
	Type                         FieldType
}

// The representation of a go struct. It has a name, a set of Field types and a Level
// to determine at what level should the struct be defined in the final source code.
type GoStruct struct {
	Name   string
	Fields []Field
	Level  int
}

func (gs GoStruct) ToStruct() string {
	var buf []string
	buf = append(buf, fmt.Sprintf("type %s struct {", gs.Name))
	for _, fld := range gs.Fields {
		var tp string
		switch fld.Type {
		case Int:
			tp = "int"
		case BigInt:
			tp = "int64"
		case Float32:
			tp = "float32"
		case Float64:
			tp = "float64"
		case String:
			tp = "string"
		case Slice:
			tp = fmt.Sprintf("[]%s", fld.TypeString)
		case Map:
			tp = fld.TypeString
		}
		buf = append(buf, fmt.Sprintf("%s %s", fld.Name, tp))
	}
	buf = append(buf, "}")
	return strings.Join(buf, "\n")
}

// Represents the FieldType in a GoStruct, an alias over int
type FieldType uint

const (
	Int = iota
	BigInt
	Float32
	Float64
	String
	Slice
	Map
)

// The decoded data represented as a struct
type DecData struct {
	mapData   map[string]interface{}
	sliceData []interface{}
}

// Convertes this instance of DecData to hold a map type value
func (dd *DecData) OfMap(mp map[string]interface{}) {
	dd.mapData = mp
}

// Converts this instance of DecData to hold a slice type value
func (dd *DecData) OfSlice(sl []interface{}) {
	dd.sliceData = sl
}

// Returns the map value from this DecData
func (dd DecData) ToMap() (map[string]interface{}, error) {
	if dd.mapData == nil {
		return nil, errors.New("MapData is not set")
	}
	return dd.mapData, nil
}

// Returns the slice value from this DecData
func (dd DecData) ToSlice() ([]interface{}, error) {
	if dd.sliceData == nil {
		return nil, errors.New("SliceData is not set")
	}
	return dd.sliceData, nil
}

// Interface type which can be decoded for parsing to go struct
type Decoder interface {
	Decode() (DecData, error)
}

// Interface type that works on annotating fields of a go struct
type Annotater interface {
	Annotate(string) string
}

// Composite Decode and Annotater interface
type DecodeAnnotater interface {
	Decoder
	Annotater
}
