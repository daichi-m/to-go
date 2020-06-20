package togo

import (
	"fmt"
	"strings"
)

// ToStruct converts the given instance of GoStruct into the actual go code for the struct.
func (gs GoStruct) ToStruct() string {
	var buf []string
	buf = append(buf, fmt.Sprintf("type %s struct {", gs.name))
	for _, fld := range gs.fields {
		var tp string

		switch fld.dataType {
		case Bool:
			tp = "bool"
		case Int:
			tp = "int"
		case Int64:
			tp = "int64"
		case Float64:
			tp = "float64"
		case String:
			tp = "string"
		case Slice:
			nest := strings.Repeat("[]", fld.sliceNesting)
			tp = fmt.Sprintf("%s%s", nest, fld.dtStruct)
		case Map:
			tp = fld.dtStruct
		}
		buf = append(buf, fmt.Sprintf("%s %s", fld.name, tp))
	}
	buf = append(buf, "}")
	return strings.Join(buf, "\n")
}
