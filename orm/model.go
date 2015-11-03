package orm

import (
	"reflect"
)

type field struct {
	Name  string
	Type  reflect.Type
	Field reflect.Value
}

type model struct {
	Type   reflect.Type
	Fields map[string]field
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

	mt := model{Fields: make(map[string]field, 0), Type: modeltype}
	for i := 0; i < value.NumField(); i++ {
		mt.Fields[modeltype.Field(i).Name] = field{modeltype.Field(i).Name, modeltype.Field(i).Type, value.Field(i)}
	}
	p.Models[modeltype] = mt
	return mt
}

func (p *model) GetValuesForSqlRowScan(cols []string) []interface{} {
	var values = make([]interface{}, len(cols))

	for index, column := range cols {
		if field, ok := p.Fields[column]; ok {
			if field.Field.Kind() == reflect.Ptr {
				values[index] = field.Field.Addr().Interface()
			} else {
				values[index] = reflect.New(reflect.PtrTo(field.Field.Type())).Interface()
			}
		} else {
			var i interface{}
			values[index] = &i
		}
	}
	return values
}

func (p *model) MapObjFromRowValues(cols []string,values []interface{}) reflect.Value {
	obj := reflect.New(p.Type).Elem()
	for index, column := range cols {
		value := values[index]
		if field, ok := p.Fields[column]; ok {
			if field.Field.Kind() == reflect.Ptr {
				obj.FieldByName(column).Set(reflect.ValueOf(value).Elem())
			} else if v := reflect.ValueOf(value).Elem().Elem(); v.IsValid() {
				obj.FieldByName(column).Set(v)
			}
		}
	}
	return obj
}
