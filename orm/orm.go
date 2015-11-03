package orm

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Orm struct {
	db      *sql.DB
	dialect Dialect
	models  *modelStruct
}

func Open(driver, source string) (*Orm, error) {
	var err error

	orm := &Orm{
		dialect: NewDialect(driver),
		models:  NewModelStruce(),
	}

	orm.db, err = sql.Open(driver, source)

	if err == nil {
		err = orm.db.Ping()
	}

	var users []User = []User{}
	test(orm, &users)
	return orm, err
}

func (p *Orm) DB() *sql.DB { return p.db }

func test(o *Orm, lst interface{}) {
	m := o.models.GetModelInfo(lst)
	fmt.Println(m)

	resultv := reflect.ValueOf(lst)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("lst argument must be a slice address")
	}
	slicev := resultv.Elem()

	rows, _ := o.db.Query("select * from [user]")
	cols, _ := rows.Columns()

	tmp_values := m.GetValuesForSqlRowScan(cols)

	for rows.Next() {
		var values = tmp_values[:]
		rows.Scan(values...)
		elem:=m.MapObjFromRowValues(cols,values)
		slicev = reflect.Append(slicev, elem)

	}
	resultv.Elem().Set(slicev.Slice(0, slicev.Cap()))
	fmt.Println(lst)
}

type User struct {
	Id  string
	Mc  string
	Age int
}
