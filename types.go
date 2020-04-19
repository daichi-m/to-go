package togo

// The struct that donates a go field inside the go struct that will be generated.
// It has a Name, a Type and optionally an Annotation to user for (un)marshalling
type field struct {
	name        string
	annotation  string
	fldType     fieldType
	fldTypeName string
	nesting     int
}

// The representation of a go struct. It has a name, a set of Field types and a Level
// to determine at what level should the struct be defined in the final source code.
type goStruct struct {
	name     string
	fields   []field
	strLevel int
}

// Represents the FieldType in a GoStruct, an alias over int
type fieldType uint

const (
	Int = iota
	Bool
	BigInt
	Float32
	Float64
	String
	Slice
	Map
)

// The decoded data represented as a struct
type decodedData struct {
	mapData   map[string]interface{}
	sliceData []interface{}
}

// Convertes this instance of DecData to hold a map type value
func (dd *decodedData) Map(mp map[string]interface{}) {
	dd.mapData = mp
}

// Converts this instance of DecData to hold a slice type value
func (dd *decodedData) Slice(sl []interface{}) {
	dd.sliceData = sl
}

// Interface type which can be decoded for parsing to go struct
type Decoder interface {
	Decode() (decodedData, error)
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
