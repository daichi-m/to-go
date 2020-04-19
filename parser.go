package togo

import (
	"errors"
	"log"
	"reflect"
)

var levelOrder map[int]goStruct

type tracker struct {
	name    string
	level   int
	nesting int
}

func Parse(dec Decoder) error {

	data, err := dec.Decode()
	if err != nil {
		log.Println("Error while decoding data", err)
		return err
	}
	var gs *goStruct

	tr := tracker{
		name:    "Document",
		level:   0,
		nesting: 0,
	}

	if data.mapData != nil {
		mp := data.mapData
		gs, err = handleMap(mp, tr)
		if err != nil {
			log.Fatal("Error while handling interface", err)
		}
	} else if data.sliceData != nil {
		sl := data.sliceData
		gs, err = handleSlice(sl, tr)
		if err != nil {
			log.Fatal("Error while handling interface", err)
		}
	}
	log.Printf("%+v %v \n", gs, err)
	str := gs.ToStruct()
	log.Printf(str)
	return nil
}

func handleMap(src map[string]interface{}, tr tracker) (*goStruct, error) {

	gs := new(goStruct)
	gs.name = tr.name

	var child *goStruct
	//var err error

	for key, val := range src {
		prmtv := false
		fld := new(field)

		tp := reflect.ValueOf(val).Kind()
		switch tp {
		case reflect.Slice:
			prmtv = false
			sl := val.([]interface{})
			ctr := tracker{
				name:    key,
				level:   tr.level + 1,
				nesting: tr.nesting + 1,
			}
			child, _ = handleSlice(sl, ctr)

			fld.fldType = Slice
			fld.fldTypeName = child.name
			fld.name = key
			fld.annotation = ""
			fld.nesting = tr.nesting
		case reflect.Map:
			prmtv = false
			mp := val.(map[string]interface{})
			ctr := tracker{
				name:    key,
				level:   tr.level + 1,
				nesting: 0,
			}
			child, _ = handleMap(mp, ctr)

			fld.fldType = Map
			fld.fldTypeName = child.name
			fld.name = key
			fld.annotation = ""
			fld.nesting = 0
		default:
			prmtv = true

		}

		if prmtv {
			fld, err := primitiveField(key, val)
			if err != nil {
				log.Fatal("Failed to convert primitive type", err)
			}
			gs.fields = append(gs.fields, *fld)
		}
	}
	return gs, nil
}

func primitiveField(key string, val interface{}) (*field, error) {

	fld := new(field)
	fld.name = key
	switch val.(type) {
	case bool:
		fld.fldType = Bool
	case int, uint, int8, uint8, int16, uint16, int32, uint32:
		fld.fldType = Int
	case int64, uint64:
		fld.fldType = BigInt
	case float32:
		fld.fldType = Float32
	case float64:
		fld.fldType = Float64
	case string:
		fld.fldType = String
	default:
		return nil, errors.New("Non primitive type")
	}
	return fld, nil
}

func handleSlice(src []interface{}, tr tracker) (*goStruct, error) {
	var gs goStruct
	for _, v := range src {
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
