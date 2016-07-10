package gobreak

import "reflect"

func Obj2M(o T) M {
	m := make(M, 0)
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		m[t.Field(i).Name] = v.Field(i).Interface()
	}
	return m
}

func M2Obj(m M, o T) {
	structValue := reflect.ValueOf(o).Elem()
	for k, v := range m {
		structFieldValue := structValue.FieldByName(k)
		if !structFieldValue.IsValid() || !structFieldValue.CanSet() {
			continue
		}
		structFieldType := structFieldValue.Type()
		val := reflect.ValueOf(v)
		if structFieldType != val.Type() {
			//TODO try convert here
			continue
		}
		structFieldValue.Set(val)
	}
}

func Ext(dist, src T) {
	m := Obj2M(src)
	M2Obj(m, dist)
}
