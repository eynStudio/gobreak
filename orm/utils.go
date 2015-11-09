package orm

import (
	"reflect"

	. "github.com/eynstudio/gobreak"
)

func Extend(dest, src T) {
	destModel :=getModelInfo(dest)
	srcModel := getModelInfo(src)
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	srcVal := reflect.Indirect(reflect.ValueOf(src))

	for k := range destModel.Fields {
		if _, ok := srcModel.Fields[k]; ok {
			val := srcVal.FieldByName(k)
			destVal.FieldByName(k).Set(val)
		}
	}
}
