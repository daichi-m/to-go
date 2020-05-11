package togo

// The decoded data represented as a struct.
// Data can be either decoded into a map[string]interface{} or a []interface{}.
type decodedData struct {
	mapData   map[string]interface{}
	sliceData []interface{}
}

// Converts this instance of decodedData to hold a map type value
func (dd *decodedData) Map(mp map[string]interface{}) {
	dd.mapData = mp
}

// Converts this instance of decodedData to hold a slice type value
func (dd *decodedData) Slice(sl []interface{}) {
	dd.sliceData = sl
}

// Interface type which can be decoded for parsing to a go struct
type Decoder interface {
	Decode() (decodedData, error)
}

// FIXME: Is this required at all?
// Interface type that works on annotating fields of a a go struct
type Annotater interface {
	Annotate(string) string
}

// Composite Decode and Annotater interface
type DecodeAnnotater interface {
	Decoder
	Annotater
}
