package orm

import (
	"database/sql"
	"fmt"

	. "github.com/eynstudio/gobreak"
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

	orm.test()
	return orm, err
}

func (p *Orm) DB() *sql.DB { return p.db }

func (p *Orm) test() {
	//	u := &User{"aaaa4", "sss", 9990}
	//	p.Update(u)

	fmt.Println(p.Count(&User{}))
	fmt.Println(p.Where("age=?", 99).Count(&User{}))
	fmt.Println(p.WhereId("aaaa").Count(&User{}))
	fmt.Printf("Has:%v\n", p.Has(&User{}, "aaaa"))

	var user User
	p.Where("age=?", 99).One(&user)
	fmt.Printf("One:%v\n", user)
	p.One(&user)
	fmt.Printf("One:%v\n", user)

	var users []User = []User{}
	p.Where("age=?", 99).All(&users)
	fmt.Println(users)
	users = []User{}
	p.All(&users)
	fmt.Println(users)

}

func (p *Orm) Where(sql string, args ...interface{}) *Scope {
	scope := NewScope(p)
	scope.Where(sql, args...)
	return scope
}
func (p *Orm) WhereId(id interface{}) *Scope {
	scope := NewScope(p)
	scope.WhereId(id)
	return scope
}
func (p *Orm) Has(data T, id interface{}) bool {
	scope := NewScope(p)
	return scope.WhereId(id).Count(data) > 0
}

func (p *Orm) Count(data T) int {
	scope := NewScope(p)
	return scope.Count(data)
}

func (p *Orm) All(data T) T {
	scope := NewScope(p)
	scope.All(data)
	return data
}

func (p *Orm) One(data T) T {
	NewScope(p).One(data)
	return data
}

func (p *Orm) Insert(data T) *Orm {
	m := p.models.GetModelInfo(data)
	builder := newSqlBuilder(&m, p.dialect)
	sql, args := builder.buildInsert(data)
	if _, err := p.db.Exec(sql, args...); err != nil {
		fmt.Println(err)
	}

	return p
}

func (p *Orm) Update(data T) *Orm {
	m := p.models.GetModelInfo(data)
	builder := newSqlBuilder(&m, p.dialect)
	sql, args := builder.buildUpdate(data)

	fmt.Println(sql, args)
	if _, err := p.db.Exec(sql, args...); err != nil {
		fmt.Println(err)
	}

	return p
}

type User struct {
	Id  string
	Mc  string
	Age int
}
