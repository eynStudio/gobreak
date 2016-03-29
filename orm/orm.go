package orm

import (
	"database/sql"

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

func MustOpen(driver, source string) *Orm {
	o, e := Open(driver, source)
	Must(e)
	return o
}

func (p *Orm) DB() *sql.DB { return p.db }

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
func (p *Orm) Query(data T, query string, args ...interface{}) T {
	NewScope(p).Query(data, query, args...)
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

func (p *Orm) UpdateFields(data T, fields []string) *Orm {
	NewScope(p).UpdateFields(data, fields)
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

func (p *Orm) Begin() (ts *TxScope, err error) {
	ts = &TxScope{}
	ts.Tx, err = p.db.Begin()
	return
}

func (p *Orm) RawCount(query string, args ...interface{}) (count int64) {
	if err := p.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0
	}
	return
}
func (p *Orm) Transact(txFunc func(*TxScope)) (err error) {
	tx, err := p.Begin()
	if err != nil {
		return
	}

	defer func() {
		err = tx.Commit()
		//		if p := recover(); p != nil {
		//			switch p := p.(type) {
		//			case error:
		//				err = p
		//			default:
		//				err = fmt.Errorf("%s", p)
		//			}
		//		}
		//		if err != nil {
		//			tx.Rollback()
		//			return
		//		}
		//		err = tx.Commit()
	}()
	txFunc(tx)
	return
}
