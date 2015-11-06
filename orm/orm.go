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
	u := &User{"aaaa4", "sss", 9990}
	p.Update(u)

	var users []User = []User{}
	p.Find(&users)
	fmt.Println(users)

	var user User
	p.First(&user)
	fmt.Println(user)


}
// Has ,if where has data,check with where, else check the data.Id
func (p *Orm) Has(data T,where ...T) bool{
	if len(where)>0 {
		
	}else{
		
	}
	return false
}
func (p *Orm) Find(data T, where ...T) *Orm {
	m := p.models.GetModelInfo(data)
	builder := newSqlBuilder(&m, p.dialect)
	rows, _ := p.db.Query(builder.buildSelect())
	m.MapRowsAsLst(rows, data)
	return p
}

func (p *Orm) First(data T, where ...T) *Orm {
	m := p.models.GetModelInfo(data)
	builder := newSqlBuilder(&m, p.dialect)
	rows, _ := p.db.Query(builder.buildSelect1())
	m.MapRowsAsObj(rows, data)
	return p
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
	
	fmt.Println(sql,args)
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
