package togo

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"
)

// JSON type structure to convert to go struct
type JSON struct {
	File string
}

// Decode this Json instance into decodedData
func (j *JSON) Decode() (DecodedData, error) {

	dd := new(DecodedData)
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
		dd.mapData = mp
	case reflect.Slice:
		sl := val.([]interface{})
		dd.sliceData = sl
	default:
		log.Println("Unknown type to decode", tp.Kind())
		return *dd, errors.New("Unknown type to decode")
	}
	return *dd, nil
}

// Annotate a string with
func (j *JSON) Annotate(string) string {
	return ""
}
