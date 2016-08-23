package orm

import (
	"database/sql"
	"encoding/json"
	"log"
	"reflect"
	"strings"

	. "github.com/eynstudio/gobreak"
)

const (
	timeFormate = "2006-01-02 15:04:05"
)

type field struct {
	Name  string
	Type  reflect.Type
	Field reflect.Value
	reflect.StructField
}

func (p field) IsJsonb() bool {
	return strings.Contains(p.StructField.Tag.Get("db"), "jsonb")
}

type model struct {
	Name   string
	Type   reflect.Type
	Fields map[string]field
	IdName string
	dbName string
}

type models struct {
	orm     *Orm
	typeMap map[reflect.Type]*model
	nameMap map[string]*model
}

func newModels(orm *Orm) *models {
	return &models{orm: orm,
		typeMap: make(map[reflect.Type]*model, 0),
		nameMap: make(map[string]*model, 0)}
}

func (p *models) get(m interface{}) *model {
	value := reflect.Indirect(reflect.ValueOf(m))
	if value.Kind() == reflect.Slice {
		value = reflect.Indirect(reflect.New(value.Type().Elem()))
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	modeltype := value.Type()

	if mt, ok := p.typeMap[modeltype]; ok {
		return mt
	}

	mt := &model{
		Name:   modeltype.Name(),
		Fields: make(map[string]field, 0),
		Type:   modeltype,
	}
	if p.orm.mapper != nil {
		mt.dbName = p.orm.mapper(mt.Name)
	} else {
		mt.dbName = mt.Name
	}
	for i := 0; i < value.NumField(); i++ {
		fs := modeltype.Field(i)
		mt.Fields[fs.Name] = field{Name: fs.Name, Type: fs.Type, StructField: fs, Field: value.Field(i)}
	}
	p.typeMap[modeltype] = mt
	p.nameMap[mt.dbName] = mt
	return mt
}

//func newModel(modelType reflect.Type) model {
//	m := model{
//		Name:   modelType.Name(),
//		Fields: make(map[string]field, 0),
//		Type:   modelType,
//	}
//	return m
//}

//func getModelInfo(val interface{}) model {
//	value := reflect.Indirect(reflect.ValueOf(val))
//	if value.Kind() == reflect.Slice {
//		value = reflect.Indirect(reflect.New(value.Type().Elem()))
//	}
//	if value.Kind() == reflect.Ptr {
//		value = value.Elem()
//	}
//	modeltype := value.Type()

//	if mt, ok := models[modeltype]; ok {
//		return mt
//	}

//	mt := newModel(modeltype)
//	for i := 0; i < value.NumField(); i++ {
//		fs := modeltype.Field(i)
//		mt.Fields[fs.Name] = field{Name: fs.Name, Type: fs.Type, StructField: fs, Field: value.Field(i)}
//	}
//	models[modeltype] = mt
//	return mt
//}

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
			if field.IsJsonb() {
				var vs []byte
				values[index] = &vs
			} else if field.Field.Kind() == reflect.Ptr {
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
			if field.IsJsonb() {
				xx := reflect.ValueOf(value).Elem().Interface().([]byte)
				fieldobj := reflect.New(field.Type).Interface()
				if err := json.Unmarshal(xx, &fieldobj); err != nil {
					log.Println(err)
				}
				log.Println(field.Type, fieldobj)
				if reflect.ValueOf(fieldobj).IsValid() {
					obj.FieldByName(column).Set(reflect.ValueOf(fieldobj).Elem())
				}
			} else if field.Field.Kind() == reflect.Ptr {
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
	var values = p.GetValuesForSqlRowScan(cols)
	rows.Scan(values...)
	elem := p.MapObjFromRowValues(cols, values)
	resultv.Elem().Set(elem)
}

func (p *model) Obj2Map(data T) map[string]interface{} {
	val := reflect.ValueOf(data)
	m := make(map[string]interface{}, len(p.Fields))
	for k, v := range p.Fields {
		if v.IsJsonb() {
			fv := val.Elem().FieldByName(k).Interface()
			xx, err := json.Marshal(fv)
			if err != nil {
				LogErr(err)
			} else {
				m[k] = xx
			}
		} else {
			m[k] = val.Elem().FieldByName(k).Interface()
		}
	}
	return m
}
