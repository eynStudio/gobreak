package orm

import (
	"reflect"

	. "github.com/eynstudio/gobreak"
)

func Extend(dest, src T) T {
	destModel :=getModelInfo(dest)
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	srcVal := reflect.Indirect(reflect.ValueOf(src))
	for k := range destModel.Fields {
		if val := srcVal.FieldByName(k);val.IsValid(){
			destVal.FieldByName(k).Set(val)
		}
	}
	return dest
}

func Map(dest,src T) T{
	return Extend(dest,src)
}