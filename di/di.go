package di

import (
	"fmt"
	"reflect"
)

var Root Container = New()

type Container interface {
	Apply(interface{}) error
	Invoke(interface{}) ([]reflect.Value, error)
	Exec(fv reflect.Value, args []reflect.Type) ([]reflect.Value, error)
	Map(interface{}) Container
	MapAs(interface{}, interface{}) Container
	Set(reflect.Type, reflect.Value) Container
	Get(reflect.Type) reflect.Value
	SetParent(Container) Container
}

type container struct {
	items  map[reflect.Type]reflect.Value
	parent Container
}

func New() Container {
	return &container{
		items: make(map[reflect.Type]reflect.Value),
	}
}

func (this *container) Invoke(f interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(f)
	inTypes := GetFuncArgs(t)
	return this.Exec(reflect.ValueOf(f), inTypes)
}

func (this *container) Exec(fv reflect.Value, args []reflect.Type) ([]reflect.Value, error) {
	vals, err := this.getVals(args)
	if err != nil {
		return nil, err
	}
	return fv.Call(vals), nil
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

func (this *container) Apply(val interface{}) error {
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

func (this *container) Map(val interface{}) Container {
	this.items[reflect.TypeOf(val)] = reflect.ValueOf(val)
	return this
}

func (this *container) MapAs(val interface{}, ifacePtr interface{}) Container {
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
