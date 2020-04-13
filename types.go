package togo

import (
	"crypto/sha1"
	"errors"
	"fmt"
)

// Interface that defines a Key method to ensure quick comparison
type Keyed interface {
	Key() string
}

// The struct that donates a go field inside the go struct that will be generated.
// It has a Name, a Type and optionally an Annotation to user for (un)marshalling
type Field struct {
	Name, Annotation string
	Type             FieldType
	HashKey          string
}

func (f *Field) Key() string {
	if len(f.HashKey) != 0 {
		return f.HashKey
	}

	s := fmt.Sprintf("%s:%d:%s", f.Name, f.Type, f.Annotation)
	sum := sha1.Sum([]byte(s))
	f.HashKey = string(sum[:])
	return f.HashKey
}

// The representation of a go struct. It has a name, a set of Field types and a Level
// to determine at what level should the struct be defined in the final source code.
type GoStruct struct {
	Name    string
	Fields  []Field
	Level   int
	HashKey string
}

func (gs *GoStruct) Key() string {
	if len(gs.HashKey) != 0 {
		return gs.HashKey
	}

	sl := make([]string, len(gs.Fields))
	for _, f := range gs.Fields {
		fk := f.Key()
		sl = append(sl, fk)
	}
	s := fmt.Sprintf("%s:%d:%v", gs.Name, gs.Level, sl)
	sum := sha1.Sum([]byte(s))
	gs.HashKey = string(sum[:])
	return gs.HashKey
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
