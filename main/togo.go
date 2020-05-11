package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daichi-m/togo"
)

func main() {

	js := new(togo.Json)

	logfile, err := os.Create("../logs/togo.log")
	if err != nil {
		fmt.Printf("Error in creating logfile, will log to stdout: %v\n", err)
	} else {
		log.SetOutput(logfile)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	wd, _ := os.Getwd()
	log.Printf("Current working directory: %s \n", wd)
	js.File = "../samples/json/bake.json"
	togo.Parse(js)
}
