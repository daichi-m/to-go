package togo

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func getLogger() *zap.Logger {
	if logger == nil {
		var err error
		config := zap.NewDevelopmentConfig()
		config.Encoding = "console"
		config.OutputPaths = []string{"./togo.log"}
		config.DisableStacktrace = false
		config.DisableCaller = false
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		logger, err = config.Build()
		if err != nil {
			fmt.Printf("logger could not be initialized, Example logger will be used")
			logger = zap.NewExample()
		}
	}
	return logger
}

// A trace interface to track the progression tree while converting to
// go struct from the generic object
type trace struct {
	name     string
	level    int
	nesting  int
	parent   *trace
	children []*trace
}

func (tt *trace) addChild(t *trace) {
	if len(tt.children) == 0 {
		tt.children = make([]*trace, 0)
	}
	tt.children = append(tt.children, t)
	t.parent = tt
}

func (tt *trace) createChild(name string, slice bool) *trace {
	t := new(trace)
	t.name = name
	t.level = tt.level + 1
	if slice {
		t.nesting = tt.nesting + 1
	} else {
		t.nesting = 0
	}
	tt.addChild(t)
	return t
}

func createRoot() *trace {
	return &trace{
		name:    "ROOT",
		level:   0,
		nesting: 0,
		parent:  nil,
	}
}

type Slicer interface {
	GetFieldDT() FieldDT
	GetStruct() IGoStruct
	GetNesting() int
	Primitive() bool
}

type sliceType struct {
	fdt     FieldDT
	strct   IGoStruct
	nesting int
}

var _ Slicer = (*sliceType)(nil)

func (st *sliceType) GetFieldDT() FieldDT {
	return st.fdt
}

func (st *sliceType) GetStruct() IGoStruct {
	return st.strct
}

func (st *sliceType) GetNesting() int {
	return st.nesting
}

func (st *sliceType) Primitive() bool {
	return st.fdt.Primitive()
}

// Parser implements the parsing logic to parse a generic data into a IGoStruct
type Parser interface {
	Parse(dec Decoder) (IGoStruct, error)
	handleMap(src map[string]interface{}, tr *trace) (IGoStruct, error)
	handleSlice(src []interface{}, tr *trace) (Slicer, error)
}

// TracedParser is a Parser that has tracer embedded in it.
type TracedParser struct {
	traceRoot *trace
	IFieldMaker
	IGoStructMaker
}

var _ Parser = (*TracedParser)(nil)

// Parse this instance of Decoder into a GoStruct.
// Returns an error in case of any error that occurs
func (tp *TracedParser) Parse(dec Decoder) (IGoStruct, error) {
	logger := getLogger().Sugar()
	defer logger.Sync()

	data, err := dec.Decode()
	if err != nil {
		logger.Fatal("Error while decoding data", zap.Error(err))
		return nil, err
	}
	logger.Debug("Decoded data", zap.Any("data", data))
	tp.traceRoot = createRoot()
	logger.Info("Tracer", zap.Any(tp.traceRoot.name, tp.traceRoot))
	var gs IGoStruct

	var nest int

	if data.mapData != nil {
		mp := data.mapData
		tr := tp.traceRoot.createChild("Document", false)
		gs, err = tp.handleMap(mp, tr)
		if err != nil {
			logger.Fatal("Error in handling map", zap.Error(err))
		}
	} else if data.sliceData != nil {
		sl := data.sliceData
		tr := tp.traceRoot.createChild("Document", true)
		slt, err := tp.handleSlice(sl, tr)
		if err != nil {
			logger.Fatal("Error in handling slice", zap.Error(err))
		}
		gs = slt.GetStruct()
	}
	logger.Debug("Final Result", zap.Any("gs", gs), zap.Int("nesting", nest), zap.Error(err))
	return gs, nil
}

// handleMap takes care of converting a map[string]interface{}
// into a GoStruct
func (tp *TracedParser) handleMap(src map[string]interface{}, tr *trace) (IGoStruct, error) {

	logger := getLogger().Sugar()
	defer logger.Sync()

	logger.Info("Tracer", zap.Any(tr.name, tr))
	gs := tp.MakeGoStruct(tr.name, nil, tr.level)

	logger.Debug("Iterate and fill up fields on GoStruct", zap.String("name", gs.Name()))
	for key, val := range src {
		dt, ok := NewFieldDT(val)
		// field, err := tp.MakeIField(key, val)
		if !ok {
			logger.Fatal("Error while converting val to data-type", zap.Any("value", val))
		}

		if dt.Primitive() {
			logger.Debug("Primitive value %s\n", zap.String("name", key))
			field := tp.MakeIField(key, "", dt, "", -1)
			gs.AddField(field)
		} else if dt == Map {
			logger.Debug("Map inside map", zap.String("key", key))
			mp := val.(map[string]interface{})
			chTr := tr.createChild(key, false)
			cgs, err := tp.handleMap(mp, chTr)
			if err != nil {
				logger.Fatal("Error in handling map", zap.Error(err))
				return nil, err
			}
			field := tp.MakeIField(key, "", dt, cgs.Name(), -1)
			gs.AddField(field)
		} else if dt == Slice {
			logger.Debug("Found a slice inside a map.", zap.String("key", key))
			sl := val.([]interface{})
			chTr := tr.createChild(key, true)
			slt, err := tp.handleSlice(sl, chTr)
			if err != nil {
				logger.Fatal("Failed converting slice to GoStruct\n", zap.Error(err))
				return nil, err
			}
			var field IField
			if slt.Primitive() {
				field = tp.MakeIField(key, "", dt, "", slt.GetNesting())
			} else {
				field = tp.MakeIField(key, "", dt, slt.GetStruct().Name(), slt.GetNesting())
			}
			gs.AddField(field)
		} else {
			logger.Info("Unknown data type, ignoring", zap.String("type", dt.GoString()))
		}
	}
	return gs, nil
}

// handleSlice takes care of converting a slice of interface{}
// into an instance of GoStruct
func (tp *TracedParser) handleSlice(src []interface{}, tr *trace) (Slicer, error) {

	logger := getLogger().Sugar()
	defer logger.Sync()
	logger.Info("Tracer", zap.Any(tr.name, tr))

	var err error

	chldGs := (IGoStruct)(nil)
	gs0 := (IGoStruct)(nil)
	dt0 := (FieldDT)(Initial)
	slType := (*sliceType)(nil)

	/*&sliceType{
		fdt:     dt,
		nesting: tr.nesting,
		strct:   nil,
	}
	*/

	for idx, val := range src {
		dt, ok := NewFieldDT(val)
		if !ok {
			logger.Fatal("Error while converting value to field", zap.Any("value", val))
			return nil, errors.New("Error converting value to field")
		}
		if idx == 0 {
			dt0 = dt
		} else if dt != dt0 {
			logger.Errorf("Different data-types found inside a list. Expected: %v, Found: %v\n",
				dt0, dt)
			return nil, fmt.Errorf("Slice not feasible. Expected %v, Found %v", dt0, dt)
		}

		if dt.Primitive() {
			if slType == nil {
				slType = &sliceType{
					fdt:     dt,
					nesting: tr.nesting,
					strct:   nil,
				}
			}
			continue
		} else if dt == Slice {
			sl := val.([]interface{})
			chtr := tr.createChild(tr.name, true)
			slTp, err := tp.handleSlice(sl, chtr)
			if err != nil {
				logger.Fatal("Could not convert the slice to GoStruct", zap.Error(err))
				return nil, err
			}

			if slTp.Primitive() {
				if slType == nil {
					slType = slTp.(*sliceType)
				}
			} else {
				chldGs = slTp.GetStruct()
			}
		} else if dt == Map {
			mp := val.(map[string]interface{})
			chtr := tr.createChild(tr.name, false)
			chldGs, err = tp.handleMap(mp, chtr)
			if err != nil {
				logger.Fatal("Could not convert the map to GoStruct", zap.Error(err))
				return nil, err
			}
		}

		if gs0 == nil {
			gs0 = chldGs
			slType = &sliceType{
				fdt:     dt0,
				nesting: tr.nesting,
				strct:   gs0,
			}
		} else {
			_, err := gs0.Grow(chldGs)
			if err != nil {
				logger.Fatal("Cannot group the existing GoStruct", zap.Error(err))
				return nil, err
			}
		}
	}
	return slType, nil
}
