package main

import (
	"log"
	"os"

	"github.com/daichi-m/togo"
)

func main() {

	js := new(togo.JSON)
	wd, _ := os.Getwd()
	log.Printf("Current working directory: %s \n", wd)
	js.File = "../samples/json/bake.json"
	fldMaker := togo.FieldMaker{}
	gsMaker := togo.GoStructMaker{}
	parser := togo.TracedParser{
		IFieldMaker:    &fldMaker,
		IGoStructMaker: &gsMaker,
	}
	parser.Parse(js)
}
