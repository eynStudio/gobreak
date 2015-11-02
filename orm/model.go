package orm

import (
	"fmt"
	"reflect"
)

type field struct {
	Name string
	Type reflect.Type
}

type model struct {
	Fields []field
}

type modelStruct struct {
	Models map[reflect.Type]model
}

func NewModelStruce() *modelStruct {
	return &modelStruct{
		Models: make(map[reflect.Type]model, 0),
	}
}

func (p *modelStruct) GetModelInfo(val interface{}) model {
	value := reflect.Indirect(reflect.ValueOf(val))
	if value.Kind() == reflect.Slice {
		value = reflect.Indirect(reflect.New(value.Type().Elem()))
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	modeltype := value.Type()

	if mt, ok := p.Models[modeltype]; ok {
		return mt
	}

	mt := model{}
	for i := 0; i < value.NumField(); i++ {
		mt.Fields = append(mt.Fields, field{modeltype.Field(i).Name, modeltype.Field(i).Type})
	}
	p.Models[modeltype] = mt
	return mt
}
