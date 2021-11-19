package utils

import (
	"errors"
	"reflect"
)

//ArrayValues returns strut values as an array
func ArrayValues(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	return values
}

//ArrayContains checks if exists in array
func ArrayContains(arr []interface{}, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

//ExistsInStrut if a value exists returns true
func ExistsInStrut(arr interface{}, toCheck string) bool {
	v := reflect.ValueOf(arr)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == toCheck {
			return true
		}
	}
	return false
}

//GetStructFieldByString returns the named field as an interface{}, also returns the Type of the interface.
func GetStructFieldByString(arr interface{}, toGet string) (interface{}, string, error) {
	v := reflect.ValueOf(arr)
	t := v.Kind().String()
	if t != "struct" {
		err := errors.New("GetStrutFieldByString(): input interface is not type Struct")
		return nil, "", err
	}
	f := reflect.Indirect(v).FieldByName(toGet)
	if !f.IsValid() {
		err := errors.New("GetStrutFieldByString(): cannot find field in struct")
		return nil, "", err
	}
	//return f.Interface(), f.Type().String(), nil
	return f.Interface(), f.Kind().String(), nil
}
