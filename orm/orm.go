package orm

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Orm struct {
	db      *sql.DB
	dialect Dialect
}

func Open(driver, source string) (*Orm, error) {
	var err error

	orm := &Orm{}
	orm.dialect = NewDialect(driver)
	orm.db, err = sql.Open(driver, source)

	if err == nil {
		err = orm.db.Ping()
	}
	
	var users []User=[]User{}
	test(orm,users)
	return orm, err
}

func (p *Orm) DB() *sql.DB { return p.db }

func test(o *Orm ,lst interface{}) {

value := reflect.Indirect(reflect.ValueOf(lst))

	if value.Kind() == reflect.Slice {
		value = reflect.Indirect(reflect.New(value.Type().Elem()))
	}
	
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
//		v2:=&value
		fmt.Println(value)
		
v := reflect.ValueOf(lst)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		
		fmt.Println("Not Struce")
				fmt.Println(v.Kind())
		return 
	}

	t := v.Type()
fmt.Println(t)

	rows, _ := o.db.Query("select * from [user]")
	cols, _ := rows.Columns()

	for rows.Next() {
		var values = make([]interface{}, len(cols))
		for i, _ := range cols {
			var a interface{}
			values[i] = &a
		}
		rows.Scan(values...)
		fmt.Println(values)
//		for _, c := range values {
//			fmt.Println(c)
//		}
	}
}

type User struct {
	Id  string
	Mc  string
	Age int
}
