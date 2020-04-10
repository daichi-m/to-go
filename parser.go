package togo

import (
	"log"
	"reflect"
)

func Parse(dec Decoder) error {

	data, err := dec.Decode()
	if err != nil {
		log.Println("Error while decoding data", err)
		return err
	}
	mp, err := data.ToMap()
	if err == nil {
		handleMap(mp)
	}
	sl, err := data.ToSlice()
	if err == nil {
		handleSlice(sl)
	}
	return nil
}

func handleMap(m map[string]interface{}) (gs GoStruct, err error) {

	var set bool
	var gs GoStruct
	var fields []Field
	for k, v := range m {
		set = false
		fld := new(Field)
		fld.Name = k
		switch v.(type) {
		case int:
			fld.Type = Int
			set = true
		case int64:
			fld.Type = BigInt
			set = true
		case string:
			fld.Type = String
			set = true
		case float32:
			fld.Type = Float32
			set = true
		case float64:
			fld.Type = Float64
			set = true
		default:
			set = false
		}

		if !set {
			tp := reflect.ValueOf(v).Kind()
			switch tp {
			case reflect.Slice:
				sl := v.([]interface{})
				mgs, err := handleSlice(sl)
				if err != nil {
					return gs, err
				}
			case reflect.Map:
				mp := v.(map[string]interface{})
				sgs, err := handleMap(mp)
				if err != nil {
					return gs, err
				}
			}
		}
		fields = append(fields, *fld)
	}
	gs.Fields = fields
	return gs, nil
}

func handleSlice(s []interface{}) (gs GoStruct, err error) {
	for i, v := range s {
		log.Printf("%d -> %v.(%T)\n", i, v, v)
	}
}
