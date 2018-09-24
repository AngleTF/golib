package fastcheck

import (
	"reflect"
)

func IsZero(v interface{}) bool {
	return v == 0 || v == "0"
}

func IsEmpty(v string) bool {
	return v == ""
}

func IsNil(v interface{}) bool{
	return v == nil
}

func IsEmptySlice(v interface{}) bool {
	var vof reflect.Value
	vof = reflect.ValueOf(v)
	return vof.Kind() == reflect.Slice && vof.Len() == 0
}