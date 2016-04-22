package di

import (
	"reflect"
)

func Ptr2Elem(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
func InterfaceOf(ifacePtr interface{}) reflect.Type {
	t := reflect.TypeOf(ifacePtr)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Interface {
		panic("value is not a pointer to an interface. (*Interface)(nil)")
	}
	return t
}

func GetFuncArgs(t reflect.Type) []reflect.Type {
	l := t.NumIn()
	in := make([]reflect.Type, l)
	for i := 0; i < l; i++ {
		in[i] = t.In(i)
	}
	return in
}

func GetMethodArgs(t reflect.Type) []reflect.Type {
	l := t.NumIn() - 1
	in := make([]reflect.Type, l)
	for i := 0; i < l; i++ {
		in[i] = t.In(i + 1)
	}
	return in
}
