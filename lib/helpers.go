package lib

import "reflect"

func SetValue(v any, field string) {
	typ := reflect.ValueOf(v).FieldByName(field).Type()
	var inputStruct reflect.Value
	if typ.Kind() == reflect.Ptr {
		inputStruct = reflect.New(typ.Elem())
	} else {
		inputStruct = reflect.New(typ).Elem()
	}

	reflect.ValueOf(&v).Elem().FieldByName(field).Set(inputStruct)
}
