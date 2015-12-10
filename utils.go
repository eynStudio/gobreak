package gobreak

import (
	"log"
	"errors"
	"reflect"
)

var (
	ErrNil         = errors.New("passed is nil")
	ErrNotSlicePtr = errors.New("passed must be slice pointer")
)

func Is(t reflect.Value, k reflect.Kind) bool { return t.Type().Kind() == k }

func IsStrT(t T) bool            { return IsStr(reflect.ValueOf(t)) }
func IsStr(t reflect.Value) bool { return Is(t, reflect.String) }

func IsPtrT(t T) bool            { return IsPtr(reflect.ValueOf(t)) }
func IsPtr(t reflect.Value) bool { return Is(t, reflect.Ptr) }

func IsSliceT(t T) bool            { return IsSlice(reflect.ValueOf(t)) }
func IsSlice(t reflect.Value) bool { return Is(t, reflect.Slice) }

func MustSlicePtr(t T) {
	if t == nil {
		panic(ErrNil)
	}

	if v := reflect.ValueOf(t); !IsPtr(v) || !IsSlice(v.Elem()) {
		panic(ErrNotSlicePtr)
	}
}


func Must(err error) (ok bool){
	if err!=nil{
		log.Fatalln(err)
	}
	return true
}