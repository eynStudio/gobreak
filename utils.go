package gobreak

import (
	"reflect"
)

func IsString(v interface{}) bool {
	return reflect.TypeOf(v).String() == "string"
}
