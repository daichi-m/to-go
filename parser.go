package togo

import (
	"errors"
	"log"
	"reflect"
)

var levelOrder map[int]goStruct

func Parse(dec Decoder) error {

	data, err := dec.Decode()
	if err != nil {
		log.Println("Error while decoding data", err)
		return err
	}
	var gs goStruct
	if data.mapData != nil {
		mp := data.mapData
		gs, err = handleMap("Document", mp, 0)
		if err != nil {
			log.Fatal("Error while handling interface", err)
		}
	} else if data.sliceData != nil {
		sl := data.sliceData
		gs, err = handleSlice("Document", sl, 0)
		if err != nil {
			log.Fatal("Error while handling interface", err)
		}
	}
	log.Printf("%+v %v \n", gs, err)
	str := gs.ToStruct()
	log.Printf(str)
	return nil
}

func handleMap(name string, m map[string]interface{}, lvl int) (goStruct, error) {

	var set bool
	var gs goStruct
	var fields []field
	for k, v := range m {
		set = false
		fld := new(field)
		fld.name = k
		switch v.(type) {
		case byte, int8, int16, int32, int:
			fld.fldType = Int
			set = true
		case int64:
			fld.fldType = BigInt
			set = true
		case string:
			fld.fldType = String
			set = true
		case float32:
			fld.fldType = Float32
			set = true
		case float64:
			fld.fldType = Float64
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
	gs.fields = fields
	gs.name = name
	gs.strLevel = lvl
	return gs, nil
}

func handleSlice(name string, s []interface{}, lvl int) (goStruct, error) {
	var gs goStruct
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
