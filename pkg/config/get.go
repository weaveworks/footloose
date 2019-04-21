package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func pathSplit(r rune) bool {
	return r == '.' || r == '[' || r == ']' || r == '"'
}

func GetValueFromConfig(stringPath string, object interface{}) (interface{}, error) {
	keyPath := strings.FieldsFunc(stringPath, pathSplit)
	// log.Println(stringPath)
	v := reflect.ValueOf(object)
	for _, key := range keyPath {
		// log.Println(key)
		keyUpper := strings.Title(key)
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			v = v.FieldByName(keyUpper)
			if v.IsValid() == false {
				return nil, fmt.Errorf("%v key does not exist", keyUpper)
			}
		} else if v.Kind() == reflect.Slice {
			index, errConv := strconv.Atoi(keyUpper)
			if errConv != nil {
				return nil, fmt.Errorf("%v is not an index", key)
			}
			v = v.Index(index)
		} else {
			return nil, fmt.Errorf("%v is neither a slice or a struct", v)
		}
	}
	return v.Interface(), nil
}
