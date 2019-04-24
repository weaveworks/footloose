package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	machinePattern   = "%d"
	maxPort          = 65535
	machineNameRegex = `^(?:m|M)achines\[[0-9]+\].(?:s|S)pec.(?:n|N)ame$`
	portRegex        = `^(?:m|M)achines\[[0-9]+\].(?:s|S)pec.(?:p|P)ortMappings\[[0-9]+\].(?:(?:h|H)ostPort|(?:c|C)ontainerPort)$`
)

// IsSetValueValid checks if value is valid for the given path
func IsSetValueValid(stringPath string, value string) (rerr error) {
	defer func() {
		if r := recover(); r != nil {
			rerr = fmt.Errorf(fmt.Sprint(r))
		}
	}()
	v := reflect.ValueOf(ClarifyArg(value))
	if v.Kind() == reflect.String {
		// check machine name
		re := regexp.MustCompile(machineNameRegex)
		if re.MatchString(stringPath) == true {
			if strings.Contains(v.Interface().(string), machinePattern) == false {
				return fmt.Errorf("Machine name is not valid, it should contain %v", machinePattern)
			}
		}
	} else if v.Kind() == reflect.Int {
		// check port value
		re := regexp.MustCompile(portRegex)
		if re.MatchString(stringPath) == true {
			if v.Interface().(int) > maxPort || v.Interface().(int) < 1 {
				return fmt.Errorf("Port cannot be higher than %v or lesset than 1", maxPort)
			}
		}
	}
	return nil
}

// ClarifyArg converts string to int or bool if possible
func ClarifyArg(v string) interface{} {
	intV, err := strconv.Atoi(v)
	if err == nil {
		return intV
	}
	boolV, err := strconv.ParseBool(v)
	if err == nil {
		return boolV
	}
	return v
}

// SetValueToConfig sets specific value to an object given a string path
func SetValueToConfig(stringPath string, object interface{}, newValue interface{}) (rerr error) {
	defer func() {
		if r := recover(); r != nil {
			rerr = fmt.Errorf(fmt.Sprint(r))
		}
	}()
	keyPath := strings.FieldsFunc(stringPath, pathSplit)
	v := reflect.ValueOf(object)
	for _, key := range keyPath {
		keyUpper := strings.Title(key)
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			v = v.FieldByName(keyUpper)
			if v.IsValid() == false {
				return fmt.Errorf("%v key does not exist", keyUpper)
			}
		} else if v.Kind() == reflect.Slice {
			index, errConv := strconv.Atoi(keyUpper)
			if errConv != nil {
				return fmt.Errorf("%v is not an index", key)
			}
			v = v.Index(index)
		} else {
			return fmt.Errorf("%v is neither a slice or a struct", v)
		}
	}
	newV := reflect.ValueOf(newValue)
	if v.Kind() == newV.Kind() {
		v.Set(newV)
	} else {
		return fmt.Errorf("%v type and %v type do not correspond", v, newV)
	}
	return nil
}
