package togo

// DecodedData is the decoded data represented as a struct.
// Data can be either decoded into a map[string]interface{} or a []interface{}.
type DecodedData struct {
	mapData   map[string]interface{}
	sliceData []interface{}
}

// Decoder is an interface type which can be used by togo classes
// to decode into a GoStruct instance
type Decoder interface {
	Decode() (DecodedData, error)
}

// Annotater is an interface type that works on
// annotating fields of a a go struct
// FIXME: Is this required at all?
type Annotater interface {
	Annotate(string) string
}

// DecodeAnnotater is a composite Decode and Annotater interface
type DecodeAnnotater interface {
	Decoder
	Annotater
}
