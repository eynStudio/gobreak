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
	u := &User{"insert1", "insert222", 9990}
	//	p.Update(u)
	p.Update(u)
	p.Where("age=?", 99).Del(&User{})
	p.Del(&User{Id: "insert1"})

	fmt.Println(p.Count(&User{}))
	fmt.Println(p.Where("age=?", 99).Count(&User{}))
	fmt.Println(p.WhereId("aaaa").Count(&User{}))
	fmt.Printf("Has:%v\n", p.HasId(&User{}, "aaaa"))

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
	return NewScope(p).Where(sql, args...)
}
func (p *Orm) WhereId(id interface{}) *Scope {
	return 	 NewScope(p).WhereId(id)
}
func (p *Orm) Has(data T) bool {
	return NewScope(p).Has(data)
}
func (p *Orm) HasId(data T, id interface{}) bool {
	return NewScope(p).WhereId(id).Count(data) > 0
}
func (p *Orm) Count(data T) int {
	return NewScope(p).Count(data)
}

func (p *Orm) All(data T) T {
	NewScope(p).All(data)
	return data
}

func (p *Orm) One(data T) T {
	NewScope(p).One(data)
	return data
}

func (p *Orm) Insert(data T) *Orm {
	NewScope(p).Insert(data)
	return p
}

func (p *Orm) Update(data T) *Orm {
	NewScope(p).Update(data)
	return p
}
func (p *Orm) Save(data T) *Orm {
	NewScope(p).Save(data)
	return p
}
func (p *Orm) Del(data T) *Orm {
	NewScope(p).Del(data)
	return p
}

type User struct {
	Id  string
	Mc  string
	Age int
}
