package main

import (
	"github.com/daichi-m/togo"
	"log"
	"os"
)

func main() {
	//fmt.Println("vim-go")
	f, err := os.Open("../samples/sample-map.json")
	if err != nil {
		log.Fatal("Error in opening file", err)
		os.Exit(1)
	}
	togo.Decode(f)
}
