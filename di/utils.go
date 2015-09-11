package di

import (
	"reflect"
)

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
