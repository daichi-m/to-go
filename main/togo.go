package main

import (
	"github.com/daichi-m/togo"
)

func main() {
	js := new(togo.Json)
	js.File = "../samples/sample-slice.json"
	togo.Parse(js)
}
