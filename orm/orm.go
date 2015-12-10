package orm

import (
	"database/sql"
	"fmt"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db"
)

type Orm struct {
	db      *sql.DB
	dialect Dialect
}

func Open(driver, source string) (*Orm, error) {
	var err error

	orm := &Orm{
		dialect: NewDialect(driver),
	}

	orm.db, err = sql.Open(driver, source)

	if err == nil {
		err = orm.db.Ping()
	}

	//	orm.test()
	return orm, err
}

func (p *Orm) DB() *sql.DB { return p.db }

func (p *Orm) test() {
	u := &User{"insert12", "insert222", 9990}
	//	p.Update(u)
	p.Save(u)
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

	var user2 User2
	Extend(&user2, u)
	fmt.Println(user2)

}

func (p *Orm) Where(sql string, args ...interface{}) *Scope {
	return NewScope(p).Where(sql, args...)
}
func (p *Orm) WhereId(id interface{}) *Scope {
	return NewScope(p).WhereId(id)
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
func (p *Orm) Page(model T, pf *db.PageFilter) *db.Paging {
	return NewScope(p).Page(model, pf)
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
func (p *Orm) SaveAs(dest T, src ...T) *Orm {
	s := NewScope(p)
	for _, m := range src {
		s.Save(Map(dest, m))
	}
	return p
}

func (p *Orm) Del(data T) *Orm {
	NewScope(p).Del(data)
	return p
}
func (p *Orm) DelAll(data T) *Orm {
	NewScope(p).DelAll(data)
	return p
}
func (p *Orm) DelId(data T, id interface{}) *Orm {
	NewScope(p).WhereId(id).Del(data)
	return p
}

func (p *Orm) Begin() (ts *TxScope) {
	ts = &TxScope{}
	ts.Tx, ts.Err = p.db.Begin()
	return
}

type User struct {
	Id  string
	Mc  string
	Age int
}

type User2 struct {
	Id  string
	Xm  string
	Age int
}
