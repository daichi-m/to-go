package togo

import "io"

type Field struct {
	Name string,
	Type string,
	Annotation string
}

type GoStruct struct {
	Name string,
	Fields []Field,
	Level int
}

type GoStructer interface {
	HandleMap(mp map[string]interface{}) (GoStruct, error)
	HandleSlice(sl []interface{}) (GoStruct, error)
	Decode(f io.Reader) (error)
	ToGoStruct() (string, error)
}
