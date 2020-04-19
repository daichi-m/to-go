package togo

import (
	"fmt"
	"strings"
)

func (gs goStruct) ToStruct() string {
	var buf []string
	buf = append(buf, fmt.Sprintf("type %s struct {", gs.name))
	for _, fld := range gs.fields {
		var tp string
		switch fld.fldType {
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
			tp = fmt.Sprintf("[]%s", fld.fldTypeName)
		case Map:
			tp = fld.fldTypeName
		}
		buf = append(buf, fmt.Sprintf("%s %s", fld.name, tp))
	}
	buf = append(buf, "}")
	return strings.Join(buf, "\n")
}
