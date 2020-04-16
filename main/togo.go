package main

import (
	"github.com/daichi-m/togo"
)

func main() {
	js := new(togo.Json)
	js.File = "../samples/sample-map.json"
	togo.Parse(js)
}
