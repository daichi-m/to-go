package togo

import (
	"errors"
	"fmt"
	"log"
)

// A tracker interface to track the progression tree while converting to
// go struct from the generic object
type tracker struct {
	name    string
	level   int
	nesting int
}

// Clone a give tracker and return a new instance of tracker
func (t tracker) clone() tracker {
	tn := tracker{
		name:    t.name,
		level:   t.level,
		nesting: t.nesting,
	}
	return tn
}

/*
var LevelOrderCache map[int][]*GoStruct
var NameStructCache map[string]*GoStruct
var trackerCache map[string]*tracker
var Logger *zap.Logger
*/

/*
func setLogger() error {

	Logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer Logger.Sync()
	return nil
}*/

// Parse this instance of Decoder into a GoStruct.
// Returns an error in case of any error that occurs
func Parse(dec Decoder) error {

	// if Logger == nil {
	// 	setLogger()
	// }
	// defer Logger.Sync()
	// LevelOrderCache = make(map[int][]*GoStruct)
	// NameStructCache = make(map[string]*GoStruct)
	// trackerCache = make(map[string]*tracker)

	data, err := dec.Decode()
	if err != nil {
		log.Fatalf("Error while decoding data: %+v\n", err)
		return err
	}
	var gs *GoStruct

	log.Printf("Decoded data from JSON: %+v\n", data)

	tr := tracker{
		name:    "Document",
		level:   0,
		nesting: 0,
	}
	var nest int

	if data.mapData != nil {
		mp := data.mapData
		gs, err = HandleMap(mp, tr)
		if err != nil {
			log.Fatalf("Error while handling interface: %+v\n", err)
		}
	} else if data.sliceData != nil {
		sl := data.sliceData
		gs, nest, err = HandleSlice(sl, tr)
		if err != nil {
			log.Fatalf("Error while handling interface: %v+\n", err)
		}
	}
	log.Printf("%+v %v %+v \n", gs, nest, err)
	str := gs.ToStruct()
	log.Printf(str)

	for name, gos := range NameStructCache {
		log.Printf("Name: %+v, Struct: %+v\n", name, gos.ToStruct())
	}
	return nil
}

// HandleMap takes care of converting a map[string]interface{}
// into a GoStruct
func HandleMap(src map[string]interface{}, tr tracker) (*GoStruct, error) {
	log.Printf("Tracking map element: %+v \n", tr)
	trackerCache[tr.name] = &tr

	gs := new(GoStruct)
	gs.Name = tr.name
	gs.Level = tr.level

	log.Printf("Iterate and fill up fields on GoStruct %+v\n", gs.Name)
	for key, val := range src {
		field, err := toField(key, val)
		if err != nil {
			log.Fatalf("Error while converting to Field: %+v\n", err)
		}
		prmtv := field.dataType.primitive()
		if prmtv == true {
			log.Printf("Primitive value, setting dtStruct and sliceNesting to defaults\n")
			field.dtStruct = ""
			field.sliceNesting = -1
			gs.AddField(field)
			log.Printf("Added field: %+v to the gostruct\n", field)
		} else if field.dataType == Map {
			log.Printf("Found a map inside a map. Key: %s \n", key)
			mp := val.(map[string]interface{})
			ctr := tracker{
				name:    key,
				nesting: -1,
				level:   tr.level + 1,
			}
			cgs, err := HandleMap(mp, ctr)
			if err != nil {
				log.Printf("Failed converting map to GoStruct due to %+v \n", err)
				return nil, err
			}
			field.dtStruct = cgs.Name
			field.sliceNesting = -1
			gs.AddField(field)
		} else if field.dataType == Slice {
			log.Printf("Found a slice inside a map. Key: %+v \n", key)
			sl := val.([]interface{})
			ctr := tracker{
				name:    key,
				nesting: 1,
				level:   tr.level,
			}
			cgs, nest, err := HandleSlice(sl, ctr)
			if err != nil {
				log.Printf("Failed converting slice to GoStruct due to %+v \n", err)
				return nil, err
			}
			field.dtStruct = cgs.Name
			field.sliceNesting = nest
			gs.AddField(field)
		} else {
			msg := fmt.Sprintf("Unknown data type found: %+v", field.dataType)
			log.Println(msg)
			return nil, errors.New(msg)
		}
	}
	err := Cache(gs)
	if err != nil {
		log.Printf("Error while caching: %+v\n", err)
		return nil, err
	}
	log.Printf("Map tracker element: %+v produced result %+v \n", tr, gs)
	return gs, nil
}

// HandleSlice takes care of converting a slice of interface{}
// into an instance of GoStruct
func HandleSlice(src []interface{}, tr tracker) (*GoStruct, int, error) {
	log.Printf("Tracker for slice: %+v \n", tr)
	trackerCache[tr.name] = &tr

	name := tr.name
	var dt0 FieldDT
	gs := new(GoStruct)
	var chgs *GoStruct

	for idx, val := range src {
		field, err := toField(name, val)
		if err != nil {
			log.Printf("Error while converting val to field: %+v\n", err)
			return nil, 0, err
		}
		if idx == 0 {
			dt0 = field.dataType
		} else if field.dataType != dt0 {
			log.Printf("Different data-types found inside a list. Expected: %+v, Found: %+v\n",
				dt0, field.dataType)
			return nil, 0, errors.New("Slice not feasible. Found different data-types")
		}

		prmtv := field.dataType.primitive()
		if prmtv {
			field.dtStruct = ""
			field.sliceNesting = tr.nesting
		} else if field.dataType == Slice {
			ctr := tracker{
				name:    tr.name,
				level:   tr.level,
				nesting: tr.nesting + 1,
			}
			sl := val.([]interface{})
			chgs, _, err = HandleSlice(sl, ctr)
			if err != nil {
				log.Printf("Could not convert the slice to GoStruct: %+v\n", err)
				return nil, 0, err
			}
		} else if field.dataType == Map {
			ctr := tracker{
				name:    tr.name,
				level:   tr.level,
				nesting: -1,
			}
			mp := val.(map[string]interface{})
			chgs, err = HandleMap(mp, ctr)
		}
		if gs == nil || gs.IsEmpty() {
			gs = chgs
		} else {
			err = gs.Grow(chgs)
			if err != nil {
				log.Printf("Cannot group the existing GoStruct due to %+v \n", err)
				return nil, 0, err
			}
		}
	}
	log.Printf("Slice tracker element: %+v produced result %+v with nesting %d \n",
		tr, gs, tr.nesting)
	return gs, tr.nesting, nil
}

// Cache the GoStruct into level order cache and name cache
func Cache(gs *GoStruct) error {
	gsl, ok := LevelOrderCache[gs.Level]
	if !ok {
		gsl = make([]*GoStruct, 8)
	}
	gsl = append(gsl, gs)
	LevelOrderCache[gs.Level] = gsl

	gsn, ok := NameStructCache[gs.Name]
	if !ok {
		NameStructCache[gs.Name] = gs
		return nil
	}
	eq := gs.Equals(gsn)
	if !eq {
		log.Printf("Found a GoStruct with same name, but the structs are not equal")
		return errors.New("Found a different GoStruct with same name")
	}
	return nil
}
