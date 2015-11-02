package orm

import (
	"database/sql"
)

type Orm struct {
	db      *sql.DB
	dialect Dialect
}

func Open(driver, source string) (*Orm, error) {
	var err error

	orm:=&Orm{}
	orm.dialect = NewDialect(driver)
	orm.db, err = sql.Open(driver, source)

	if err == nil {
		err = orm.db.Ping()
	}

	return orm, err
}

func (p *Orm) DB() *sql.DB { return p.db }
