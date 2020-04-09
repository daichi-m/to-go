package togo

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"reflect"
)

func Decode(f io.Reader) (int, error) {
	dec := json.NewDecoder(f)
	var value interface{}
	err := dec.Decode(&value)
	if err != nil {
		log.Print("Could not decode due to", err)
		return 0, err
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice:
		sl, ok := value.([]interface{})
		if !ok {
			log.Println("Not a Slice")
			return 0, errors.New("Not a Slice")
		}
		log.Printf("%T %v \n", sl, sl)
		HandleSlice(sl)
		return 1, nil
	case reflect.Map:
		mp, ok := value.(map[string]interface{})
		if !ok {
			log.Println("Not a Map")
			return 0, errors.New("Not a Map")
		}
		log.Printf("%T %v \n", mp, mp)
		HandleMap(mp)
		return 2, nil
	}
	return 0, errors.New("Unknown type in reflect")
}

func HandleMap(m map[string]interface{}) error {
	for k, v := range m {
		log.Printf("%s -> %v.(%T)\n", k, v, v)
	}
	return nil
}

func HandleSlice(s []interface{}) error {
	for i, v := range s {
		log.Printf("%d -> %v.(%T)\n", i, v, v)
	}
	return nil
}
