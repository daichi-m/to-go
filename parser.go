package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

func main() {
	fmt.Println("vim-go")
	f, err := os.Open("sample.json")
	if err != nil {
		log.Fatal("Error in opening file", err)
		os.Exit(1)
	}
	dec := json.NewDecoder(f)
	var v interface{}
	for dec.More() {
		err = dec.Decode(&v)
		if err != nil {
			log.Print("Error in Decode call of Decoder", err)
		}
		switch reflect.ValueOf(v).Kind() {
		case reflect.Slice:
			fmt.Println("Slice!!")
		case reflect.Map:
			fmt.Println("Map!!")
		}
		t := reflect.TypeOf(v)
		fmt.Println(t)
		fmt.Println(v)
	}
}
