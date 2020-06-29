package togo

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"go.uber.org/zap"
)

// JSON type structure to convert to go struct
type JSON struct {
	File string
}

// Decode this Json instance into decodedData
func (j *JSON) Decode() (DecodedData, error) {

	logger := getLogger().Sugar()
	defer logger.Sync()

	dd := new(DecodedData)
	f, err := os.Open(j.File)
	if err != nil {
		logger.Fatal("Error while reading file", zap.Error(err))
		return *dd, err
	}
	var val interface{}
	dec := json.NewDecoder(f)
	err = dec.Decode(&val)
	if err != nil {
		logger.Fatal("Error while decoding", zap.Error(err))
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
		logger.Fatal("Unknown type to decode", zap.Any("type", tp.Kind()))
		return *dd, errors.New("Unknown type to decode")
	}
	return *dd, nil
}

// Annotate a string with
func (j *JSON) Annotate(string) string {
	return ""
}
