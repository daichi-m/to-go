package togo

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"
)

type Json struct {
	File string
}

func (j *Json) Decode() (DecData, error) {

	dd := new(DecData)
	f, err := os.Open(j.File)
	if err != nil {
		log.Println("Error while reading file", err)
		return *dd, err
	}
	var val interface{}
	dec := json.NewDecoder(f)
	err = dec.Decode(&val)
	if err != nil {
		log.Println("Error while decoding", err)
		return *dd, err
	}
	tp := reflect.ValueOf(val)
	switch tp.Kind() {
	case reflect.Map:
		mp := val.(map[string]interface{})
		dd.OfMap(mp)
	case reflect.Slice:
		sl := val.([]interface{})
		dd.OfSlice(sl)
	default:
		log.Println("Unknown type to decode", tp.Kind())
		return *dd, errors.New("Unknown type to decode")
	}
	return *dd, nil
}
