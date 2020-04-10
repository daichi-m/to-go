package togo

/* The struct that donates a go field inside the go struct that will be generated.
 * It has a Name, a Type and optionally an Annotation to user for (un)marshalling
 */
type Field struct {
	Name, Type, Annotation string
}

/* The representation of a go struct. It has a name, a set of Field types and a Level
 * to determine at what level should the struct be defined in the final source code.
 */
type GoStruct struct {
	Name   string
	Fields []Field
	Level  int
}

/* The decoded data represented as a struct */
type DecData struct {
	MapData   map[string]interface{}
	SliceData []interface{}
}

/* Interface type that defines decodability for the parser to process */
type Decoder interface {
	Decode() (DecData, error)
}
