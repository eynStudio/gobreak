package di

import (
	"reflect"
	"testing"
)

type Test1 struct {
	Dep1 string `di:"-"`
	Dep2 int    `di`
}

func Test_DiGet(t *testing.T) {
	d := New().Map("hi")
	if !d.Get(reflect.TypeOf("string")).IsValid() {
		t.Fatal()
	}
}

func Test_DiSetParent(t *testing.T) {
	p := New().Map("hi")
	d := New().SetParent(p)
	if !d.Get(reflect.TypeOf("string")).IsValid() {
		t.Fatal()
	}
}

func Test_DiApply(t *testing.T) {
	d := New().Map("hi").Map(11)
	t1 := Test1{}
	err := d.Apply(&t1)

	if err != nil {
		t.Error(err)
	}
	if t1.Dep1 != "hi" {
		t.Fatal()
	}
	if t1.Dep2 != 11 {
		t.Fatal()
	}
}
