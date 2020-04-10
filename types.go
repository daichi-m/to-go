package togo

import "errors"

/* The struct that donates a go field inside the go struct that will be generated.
 * It has a Name, a Type and optionally an Annotation to user for (un)marshalling
 */
type Field struct {
	Name, Annotation string
	Type             FieldType
}

/* The representation of a go struct. It has a name, a set of Field types and a Level
 * to determine at what level should the struct be defined in the final source code.
 */
type GoStruct struct {
	Name   string
	Fields []Field
	Level  int
}

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

/* The decoded data represented as a struct */
type DecData struct {
	mapData   map[string]interface{}
	sliceData []interface{}
}

func (dd *DecData) OfMap(mp map[string]interface{}) {
	dd.mapData = mp
}

func (dd *DecData) OfSlice(sl []interface{}) {
	dd.sliceData = sl
}

func (dd DecData) ToMap() (map[string]interface{}, error) {
	if dd.mapData == nil {
		return nil, errors.New("MapData is not set")
	}
	return dd.mapData, nil
}

func (dd DecData) ToSlice() ([]interface{}, error) {
	if dd.sliceData == nil {
		return nil, errors.New("SliceData is not set")
	}
	return dd.sliceData, nil
}

/* Interface type that defines decodability for the parser to process */
type Decoder interface {
	Decode() (DecData, error)
	Annotation() string
}
