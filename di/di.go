package di

import (
	"fmt"
	"log"
	"reflect"

	. "github.com/eynstudio/gobreak"
)

var Root Container = New()

type Container interface {
	Apply(T) error
	Invoke(T) ([]reflect.Value, error)
	Exec(fv reflect.Value, args []reflect.Type) ([]reflect.Value, error)
	ApplyAndMap(T) error
	ApplyAndMapAs(T, T) error
	Map(T) Container
	MapAs(T, T) Container
	Set(reflect.Type, reflect.Value) Container
	Get(reflect.Type) reflect.Value
	SetParent(Container) Container
	ShowItems()
}

func Invoke(f T) ([]reflect.Value, error)               { return Root.Invoke(f) }
func Apply(val T) error                                 { return Root.Apply(val) }
func ApplyAndMap(val T) error                           { return Root.ApplyAndMap(val) }
func ApplyAndMapAs(val T, ifacePtr T) error             { return Root.ApplyAndMapAs(val, ifacePtr) }
func Map(val T) Container                               { return Root.Map(val) }
func MapAs(val T, ifacePtr T) Container                 { return Root.MapAs(val, ifacePtr) }
func Set(typ reflect.Type, val reflect.Value) Container { return Root.Set(typ, val) }
func Get(t reflect.Type) reflect.Value                  { return Root.Get(t) }

func Exec(fv reflect.Value, args []reflect.Type) ([]reflect.Value, error) { return Root.Exec(fv, args) }

type container struct {
	items  map[reflect.Type]reflect.Value
	parent Container
}

func New() Container { return &container{items: make(map[reflect.Type]reflect.Value)} }

func (this *container) Invoke(f T) ([]reflect.Value, error) {
	t := reflect.TypeOf(f)
	inTypes := GetFuncArgs(t)
	return this.Exec(reflect.ValueOf(f), inTypes)
}

func (this *container) Exec(fv reflect.Value, args []reflect.Type) ([]reflect.Value, error) {
	if vals, err := this.getVals(args); err == nil {
		return fv.Call(vals), nil
	} else {
		return nil, err
	}
}

func (this *container) getVals(types []reflect.Type) ([]reflect.Value, error) {
	l := len(types)
	vals := make([]reflect.Value, l)
	for i := 0; i < l; i++ {
		val := this.Get(types[i])
		if !val.IsValid() {
			return nil, fmt.Errorf("Value not found for type %v", types[i])
		}
		vals[i] = val
	}
	return vals, nil
}

func (this *container) Apply(val T) error {
	v := reflect.ValueOf(val)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		structField := t.Field(i)
		if f.CanSet() && (structField.Tag == "di" || structField.Tag.Get("di") != "") {
			ft := f.Type()
			v := this.Get(ft)
			if !v.IsValid() {
				return fmt.Errorf("Value not found for type %v", ft)
			}
			f.Set(v)
		}
	}
	return nil
}
func (this *container) ApplyAndMap(val T) error { return this.Map(val).Apply(val) }

func (this *container) ApplyAndMapAs(val T, ifacePtr T) error {
	return this.MapAs(val, ifacePtr).Apply(val)
}

func (this *container) Map(val T) Container {
	this.items[reflect.TypeOf(val)] = reflect.ValueOf(val)
	return this
}

func (this *container) MapAs(val T, ifacePtr T) Container {
	t := InterfaceOf(ifacePtr)
	this.items[t] = reflect.ValueOf(val)
	return this
}

func (this *container) Set(typ reflect.Type, val reflect.Value) Container {
	this.items[typ] = val
	return this
}

func (this *container) Get(t reflect.Type) reflect.Value {
	val := this.items[t]
	if val.IsValid() {
		return val
	}

	if t.Kind() == reflect.Interface {
		for k, v := range this.items {
			if k.Implements(t) {
				val = v
				break
			}
		}
	}

	if !val.IsValid() && this.parent != nil {
		val = this.parent.Get(t)
	}

	return val
}

func (this *container) SetParent(parent Container) Container {
	this.parent = parent
	return this
}

func (this *container) ShowItems() {
	log.Println(this.items)
}
