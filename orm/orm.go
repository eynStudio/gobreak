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
	
	var users []User = []User{}
	p.Find(&users)
	fmt.Println(users)
	
	var user User
	p.First(&user)
	fmt.Println(user)
	
}
func (p *Orm) Find(out T, where ...T) *Orm {
	m := p.models.GetModelInfo(out)

	rows, _ := p.db.Query("select * from [user]")

	m.MapRowsAsLst(rows, out)
	return p
}

func (p *Orm) First(out T, where ...T) *Orm {
	m := p.models.GetModelInfo(out)

	rows, _ := p.db.Query("select top 1 * from [user]")

	m.MapRowsAsObj(rows, out)
	return p
}

type User struct {
	Id  string
	Mc  string
	Age int
}
