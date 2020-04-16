package togo

import (
	"errors"
	"log"
	"reflect"
)

var structCache map[string][]*GoStruct

func Parse(da DecodeAnnotater) error {

	data, err := da.Decode()
	if err != nil {
		log.Println("Error while decoding data", err)
		return err
	}
	var gs GoStruct
	if data.mapData != nil {
		mp := data.mapData
		gs, err = handleMap("Document", mp, 0)
	} else if data.sliceData != nil {
		sl := data.sliceData
		gs, err = handleSlice("Document", sl, 0)
	}
	log.Printf("%+v %v \n", gs, err)
	str := gs.ToStruct()
	log.Printf(str)
	return nil
}

func handleMap(name string, m map[string]interface{}, lvl int) (GoStruct, error) {

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
				mgs, err := handleSlice(k, sl, lvl+1)
				if err != nil {
					return gs, err
				}
				log.Printf("%+v \n", mgs)
			case reflect.Map:
				mp := v.(map[string]interface{})
				sgs, err := handleMap(k, mp, lvl+1)
				if err != nil {
					return gs, err
				}
				log.Printf("%+v \n", sgs)
			}
		}
		fields = append(fields, *fld)
	}
	gs.Fields = fields
	gs.Name = name
	gs.Level = lvl
	return gs, nil
}

func handleSlice(name string, s []interface{}, lvl int) (GoStruct, error) {
	var gs GoStruct
	for _, v := range s {
		tp := reflect.ValueOf(v).Kind()
		switch tp {
		case reflect.Map:
			mp := v.(map[string]interface{})
			return handleMap(name, mp, lvl)
		case reflect.Slice:
			sl := v.([]interface{})
			return handleSlice(name, sl, lvl)
		default:
			log.Printf("Unknown type of slice")
		}
	}
	return gs, errors.New("The Slice contains non-map, non-slice data. Slice must contain map or slice")
}
