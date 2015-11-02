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
	test(orm, users)
	return orm, err
}

func (p *Orm) DB() *sql.DB { return p.db }

func test(o *Orm, lst interface{}) {

	m := o.models.GetModelInfo(lst)

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
