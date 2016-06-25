package orm

import (
	"database/sql"
	"reflect"
	//	"time"

	. "github.com/eynstudio/gobreak"
)

const (
	timeFormate = "2006-01-02 15:04:05"
)

var models = map[reflect.Type]model{}

type field struct {
	Name  string
	Type  reflect.Type
	Field reflect.Value
}

type model struct {
	Name   string
	Type   reflect.Type
	Fields map[string]field
	IdName string
}

func newModel(modelType reflect.Type) model {
	m := model{
		Name:   modelType.Name(),
		Fields: make(map[string]field, 0),
		Type:   modelType,
	}
	return m
}

func getModelInfo(val interface{}) model {
	value := reflect.Indirect(reflect.ValueOf(val))
	if value.Kind() == reflect.Slice {
		value = reflect.Indirect(reflect.New(value.Type().Elem()))
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	modeltype := value.Type()

	if mt, ok := models[modeltype]; ok {
		return mt
	}

	mt := newModel(modeltype)
	for i := 0; i < value.NumField(); i++ {
		mt.Fields[modeltype.Field(i).Name] = field{modeltype.Field(i).Name, modeltype.Field(i).Type, value.Field(i)}

	}
	models[modeltype] = mt
	return mt
}

func (p *model) Id() string {
	if len(p.IdName) > 0 {
		return p.IdName
	}

	if _, ok := p.Fields["Id"]; ok {
		p.IdName = "Id"
		return p.IdName
	}
	for k, v := range p.Fields {
		if v.Type.Name() == "GUID" {
			p.IdName = k
			return p.IdName
		}
	}
	panic("Can not find Id")
}
func (p *model) IdVal(obj T) interface{} {
	id := p.Id()
	val := reflect.ValueOf(obj).Elem()
	return val.FieldByName(id).Interface()
}

func (p *model) GetValuesForSqlRowScan(cols []string) []interface{} {
	var values = make([]interface{}, len(cols))

	for index, column := range cols {
		if field, ok := p.Fields[column]; ok {
			if field.Field.Kind() == reflect.Ptr {
				values[index] = field.Field.Addr().Interface()
			} else if field.Field.Kind() == reflect.String {
				values[index] = reflect.New(reflect.PtrTo(reflect.TypeOf(""))).Interface()
			} else if field.Field.Kind() == reflect.Struct {
				switch field.Field.Type().String() {
				//				case "time.Time":
				//					values[index] = reflect.New(reflect.PtrTo(reflect.TypeOf(""))).Interface()
				default:
					values[index] = reflect.New(reflect.PtrTo(field.Field.Type())).Interface()
				}
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

func (p *model) MapObjFromRowValues(cols []string, values []interface{}) reflect.Value {
	obj := reflect.New(p.Type).Elem()
	for index, column := range cols {
		value := values[index]
		if field, ok := p.Fields[column]; ok {
			if field.Field.Kind() == reflect.Ptr {
				obj.FieldByName(column).Set(reflect.ValueOf(value).Elem())
			} else if field.Field.Kind() == reflect.String {
				if v := reflect.ValueOf(value).Elem().Elem(); v.IsValid() {
					obj.FieldByName(column).SetString(v.Interface().(string))
				}
				//			} else if field.Field.Kind() == reflect.Struct {
				//				switch field.Field.Type().String() {
				//				case "time.Time":
				//					v := reflect.ValueOf(value).Elem().Elem()
				//					t, _ := time.Parse(timeFormate, v.Interface().(string))
				//					obj.FieldByName(column).Set(reflect.ValueOf(t))
				//				}
			} else if v := reflect.ValueOf(value).Elem().Elem(); v.IsValid() {
				obj.FieldByName(column).Set(v)
			}
		}
	}
	return obj
}

func (p *model) MapRowsAsLst(rows *sql.Rows, out T) {
	resultv := reflect.ValueOf(out)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("out argument must be a slice address")
	}

	slicev := resultv.Elem()
	cols, _ := rows.Columns()
	tmp_values := p.GetValuesForSqlRowScan(cols)

	for rows.Next() {
		var values = tmp_values[:]
		rows.Scan(values...)
		elem := p.MapObjFromRowValues(cols, values)
		slicev = reflect.Append(slicev, elem)
	}
	resultv.Elem().Set(slicev.Slice(0, slicev.Len()))
}

func (p *model) MapRowsAsObj(rows *sql.Rows, out T) {
	resultv := reflect.ValueOf(out)
	cols, _ := rows.Columns()
	if rows.Next() {
		var values = p.GetValuesForSqlRowScan(cols)
		rows.Scan(values...)
		elem := p.MapObjFromRowValues(cols, values)
		resultv.Elem().Set(elem)
	}
}

func (p *model) Obj2Map(data T) map[string]interface{} {
	val := reflect.ValueOf(data)
	m := make(map[string]interface{}, len(p.Fields))
	for k := range p.Fields {
		m[k] = val.Elem().FieldByName(k).Interface()
	}
	return m
}
