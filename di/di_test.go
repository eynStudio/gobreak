package di

import (
	"reflect"
	"testing"
)

type Test1 struct {
	Dep1 string `di:"-"`
	Dep2 int    `di`
}

func Test_Get(t *testing.T) {
	d := New()

	d.Map("hi")
	d.Map(13)

	t.Log(d.Get(reflect.TypeOf("string")).IsValid())
	t.Log(d.Get(reflect.TypeOf(11)).IsValid())

}
