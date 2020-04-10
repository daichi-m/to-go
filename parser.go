package togo

import (
	"log"
)

func Parse(dec Decoder) error {

	data, err := dec.Decode()
	if err != nil {
		log.Println("Error while decoding data", err)
		return err
	}
	if data.MapData != nil {
		HandleMap(data.MapData)
	} else {
		HandleSlice(data.SliceData)
	}
	return nil
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
