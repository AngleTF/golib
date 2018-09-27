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

func SliceLenEqual(args ...interface{}) bool {
	var vof reflect.Value
	var flag bool = true
	var length int
	for k, v := range args{
		vof = reflect.ValueOf(v)
		if vof.Kind() != reflect.Slice {
			flag = false
		}
		if IsZero(k){
			length = vof.Len()
		}
		if length != vof.Len(){
			flag = false
		}
	}
	return flag
}