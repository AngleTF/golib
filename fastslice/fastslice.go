package fastslice

import "reflect"

func Pop(v interface{}) (reflect.Value, bool) {
	var vof, lastSlice reflect.Value
	vof = reflect.ValueOf(v)
	if vof.Kind() != reflect.Ptr{
		return reflect.ValueOf(nil), false
	}
	vof = vof.Elem()
	if !vof.CanSet() && vof.Kind() != reflect.Slice{
		return reflect.ValueOf(nil), false
	}

	if vof.Len() == 0{
		return reflect.ValueOf(nil), false
	}

	lastSlice = vof.Slice(vof.Len() - 1, vof.Len())
	vof.Set(vof.Slice(0, vof.Len() - 1))
	return lastSlice.Index(0), true
}
