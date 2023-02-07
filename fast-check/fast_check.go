package fastCheck

import (
	"fmt"
	"reflect"
)

func IsZero(v interface{}) bool {
	return fmt.Sprint(v) == "0"
}

func IsEmpty(v string) bool {
	return v == ""
}

func IsNil(v interface{}) bool {
	return v == nil
}

func IsEmptySlice(v interface{}) bool {
	var vof reflect.Value
	vof = reflect.ValueOf(v)
	return vof.Kind() == reflect.Slice && vof.Len() == 0
}

func SliceLenEqual(args ...interface{}) bool {
	var vof reflect.Value
	var flag bool = true
	var length int
	for k, v := range args {
		vof = reflect.ValueOf(v)
		if vof.Kind() != reflect.Slice {
			flag = false
		}
		if IsZero(k) {
			length = vof.Len()
		}
		if length != vof.Len() {
			flag = false
		}
	}
	return flag
}

func InSlice(arr interface{}, in interface{}) bool {

	var (
		rv reflect.Value
	)

	rv = reflect.ValueOf(arr)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		for i := 0; i < rv.Len(); i++ {
			if rv.Index(i).Interface() == in {
				return true
			}
		}
	}

	return false
}
