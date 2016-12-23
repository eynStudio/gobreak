package orm

import (
	"database/sql"
	"fmt"
	"strings"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db"
	"github.com/eynstudio/gobreak/db/filter"
	"github.com/eynstudio/gobreak/db/meta"
)

type Orm struct {
	db      *sql.DB
	dialect Dialect
	mapper  MapperFn
	models  *models
}

func Open(driver, source string) (*Orm, error) {
	var err error

	orm := &Orm{dialect: NewDialect(driver)}
	orm.models = newModels(orm)
	orm.db, err = sql.Open(driver, source)
	orm.mapper = SameMapper
	if err == nil {
		err = orm.db.Ping()
	}
	return orm, err
}

func MustOpen(driver, source string) *Orm {
	o, e := Open(driver, source)
	Must(e)
	return o
}

func (p *Orm) getBuilder(s *Scope) Ibuilder {
	switch p.dialect.Driver() {
	case "mysql":
		return &mysqlBuilder{builder: newBuilder(s)}
	case "postgres":
		return &pgBuilder{builder: newBuilder(s)}
	case "oci8":
		return &oci8Builder{builder: newBuilder(s)}
	case "mssql":
		return &mssqlBuilder{builder: newBuilder(s)}
	}
	return nil
}
func (p *Orm) SetMapper(f MapperFn) *Orm {
	p.mapper = f
	return p
}
func (p *Orm) Mapper() MapperFn {
	if p.mapper == nil {
		return SameMapper
	}
	return p.mapper
}
func (p *Orm) DB() *sql.DB            { return p.db }
func (p *Orm) LoadMeta() *meta.MetaDb { return p.dialect.LoadMeta(p.db) }
func (p *Orm) Order(args ...string) *Scope {
	return NewScope(p).Order(args...)
}

func (p *Orm) Where(sql string, args ...interface{}) *Scope {
	return NewScope(p).Where(sql, args...)
}
func (p *Orm) WhereId(id interface{}) *Scope { return NewScope(p).WhereId(id) }

func (p *Orm) Has(data T) (bool, error) {
	s := NewScope(p)
	b := s.Has(data)
	return b, s.Err
}

func (p *Orm) HasId(data T, id interface{}) (bool, error) {
	s := NewScope(p)
	c := s.WhereId(id).Count(data)
	return c > 0, s.Err
}

func (p *Orm) Count(data T) (int, error) {
	s := NewScope(p)
	c := s.Count(data)
	return c, s.Err
}

func (p *Orm) All(data T) error { return NewScope(p).All(data).Err }
func (p *Orm) AllJson(lst T) error {
	return NewScope(p).AllJson(lst).Err
}

func (p *Orm) Query(data T, query string, args ...interface{}) error {
	return NewScope(p).Query(data, query, args...).Err
}

func (p *Orm) Page(model T, pf *filter.PageFilter) (*db.Paging, error) {
	s := NewScope(p)
	pp := s.Page(model, pf)
	return pp, s.Err
}

func (p *Orm) PageByOrder(model T, order string, pf *filter.PageFilter) (*db.Paging, error) {
	s := NewScope(p)
	pp := s.PageByOrder(model, order, pf)
	return pp, s.Err
}

func (p *Orm) Select(str string) *Scope { return NewScope(p).Select(str) }
func (p *Orm) From(name string) *Scope  { return NewScope(p).From(name) }

//func (p *Orm) One(data T) error                 { return NewScope(p).One(data).Err }
func (p *Orm) Insert(data T) error              { return NewScope(p).Insert(data).Err }
func (p *Orm) Update(data T) error              { return NewScope(p).Update(data).Err }
func (p *Orm) Save(data T) error                { return NewScope(p).Save(data).Err }
func (p *Orm) SaveTo(name string, data T) error { return NewScope(p).SaveTo(name, data).Err }
func (p *Orm) Exec(sql string, args ...interface{}) error {
	var sa db.SqlArgs
	sa.Sql = sql
	sa.Args = args
	s := NewScope(p).exec(sa)
	return s.Err
}
func (p *Orm) SaveJson(id GUID, data T) error { return NewScope(p).SaveJson(id, data).Err }
func (p *Orm) SaveJsonTo(name string, id GUID, data T) error {
	return NewScope(p).From(name).SaveJson(id, data).Err
}

//func (p *Orm) GetJson(id GUID, data T) error { return NewScope(p).GetJson(data).Err }
//func (p *Orm) GetJsonFrom(name string, id GUID, data T) error {
//	return NewScope(p).From(name).GetJson(id, data).Err
//}

func (p *Orm) UpdateFields(data T, fields []string) error {
	return NewScope(p).UpdateFields(data, fields).Err
}

//!! this must use Tx !!
//func (p *Orm) SaveAs(dest T, src ...T) error {
//	s := NewScope(p)
//	for _, m := range src {
//		s.NoErrExec(func() {
//			s.Save(Ext(dest, m))
//		})
//	}
//	return s.Err
//}

func (p *Orm) Del(data T) error                   { return NewScope(p).Del(data).Err }
func (p *Orm) DelAll(data T) error                { return NewScope(p).DelAll(data).Err }
func (p *Orm) DelId(data T, id interface{}) error { return NewScope(p).WhereId(id).Del(data).Err }

func (p *Orm) Begin() (ts *TxScope) {
	ts = NewTxScope(p)
	ts.Tx, ts.Err = p.db.Begin()
	return
}

func (p *Orm) RawCount(query string, args ...interface{}) (count int64, err error) {
	query = p.convParams(query)
	//	log.Println(query)
	err = p.db.QueryRow(query, convertArgs2(args)...).Scan(&count)
	return
}

func (p *Orm) convParams(sql string) (str string) {
	if p.dialect.Driver() != "postgres" {
		return sql
	}
	parts := strings.Split(sql, "?")
	l := len(parts)
	for i, c := range parts {
		if i < l-1 {
			str += c + fmt.Sprintf("%s%v", "$", i+1)
		}
	}
	str += parts[l-1]
	return
}

func (p *Orm) Transact(txFunc func(*TxScope)) (err error) {
	tx := p.Begin()
	if tx.IsErr() {
		return tx.Err
	}

	defer func() {
		err = tx.Commit()
	}()

	txFunc(tx)
	return
}
