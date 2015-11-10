package orm

import (
	"fmt"
	"strings"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db"
)

type Scope struct {
	orm       *Orm
	model     *model
	haswhere  bool
	where     map[string]interface{}
	whereid   interface{}
	wheresql  string
	whereargs []interface{}
	offset    int
	limit     int
	hasLimit  bool
}

func NewScope(orm *Orm) *Scope {
	return &Scope{orm: orm}
}

func (p *Scope) Where(sql string, args ...interface{}) *Scope {
	p.haswhere = true
	p.wheresql = sql
	p.whereargs = args
	return p
}

func (p *Scope) WhereId(id interface{}) *Scope {
	p.haswhere = true
	p.whereid = id
	return p
}

func (p *Scope) Count(model T) int {
	p.checkModel(model)
	id := p.model.Id()
	wsql, args := p.buildWhere()
	sql := fmt.Sprintf("SELECT COUNT(%v) from %s %v", p.quote(id), p.quote(p.model.Name), wsql)
	row := p.orm.db.QueryRow(sql, convertArgs(args)...)
	count := 0
	row.Scan(&count)
	return count
}
func (p *Scope) Has(model T) bool {
	p.checkModel(model)
	id := p.model.Id()
	idval := p.model.IdVal(model)
	sql := fmt.Sprintf("SELECT COUNT(%v) from %s WHERE %v=?", p.quote(id), p.quote(p.model.Name), id)
	row := p.orm.db.QueryRow(sql, convertArgs(idval)...)
	count := 0
	if err := row.Scan(&count); err != nil {
		fmt.Println("Has:", sql)
		fmt.Println(err)
	}
	return count > 0
}
func (p *Scope) One(model T) (has bool) {
	p.checkModel(model)
	wsql, args := p.buildWhere()
	sql := fmt.Sprintf("SELECT TOP 1 * from %s %v", p.quote(p.model.Name), wsql)
	rows, _ := p.orm.db.Query(sql, convertArgs(args)...)
	p.model.MapRowsAsObj(rows, model)
	return true
}

func (p *Scope) All(model T) {
	p.checkModel(model)
	wsql, args := p.buildWhere()
	sql := fmt.Sprintf("SELECT * from %s %v", p.quote(p.model.Name), wsql)
	rows, _ := p.orm.db.Query(sql,convertArgs(args)...)
	p.model.MapRowsAsLst(rows, model)
}

func (p *Scope) Limit(offset, limit int) *Scope {
	p.offset = offset
	p.limit = limit
	p.hasLimit = true
	return p
}
func (p *Scope) Page(model T) *db.Paging {
	paging := &db.Paging{}

	return paging
}
func (p *Scope) Save(model T) *Scope {
	p.checkModel(model)
	if p.Has(model) {
		p.Update(model)
	} else {
		p.Insert(model)
	}
	return p
}

func (p *Scope) Insert(model T) *Scope {
	p.checkModel(model)
	sql, args := p.buildInsert(model)
	p.exec(sql, args...)
	return p
}

func (p *Scope) Update(model T) *Scope {
	p.checkModel(model)
	p.setWhereIdIfNoWhere(model)
	sql, args := p.buildUpdate(model)
	p.exec(sql, args...)
	return p
}

func (p *Scope) Del(model T) *Scope {
	p.checkModel(model)
	p.setWhereIdIfNoWhere(model)
	wsql, args := p.buildWhere()
	sql := fmt.Sprintf("DELETE from %s %v", p.quote(p.model.Name), wsql)
	p.exec(sql, args...)
	return p
}

func (p *Scope) buildWhere() (string, []interface{}) {
	if !p.haswhere {
		return "", nil
	} else if p.whereid != nil {
		return fmt.Sprintf("WHERE (%v=?)", p.quote(p.model.Id())), []interface{}{p.whereid}
	} else if len(p.wheresql) > 0 {
		return "WHERE " + p.wheresql, p.whereargs
	} else if len(p.where) > 0 {

	}

	return "", nil
}

func (p *Scope) buildInsert(obj T) (string, []interface{}) {
	var cols []string
	var params []string
	var args []interface{}
	m := p.model.Obj2Map(obj)
	for k, v := range m {
		cols = append(cols, p.quote(k))
		params = append(params, "?")
		args = append(args, v)
	}
	sql := fmt.Sprintf("insert into %s (%v) values (%v)",
		p.quote(p.model.Name),
		strings.Join(cols, ","),
		strings.Join(params, ","),
	)
	return sql, args
}

func (p *Scope) buildUpdate(obj T) (string, []interface{}) {
	wsql, wargs := p.buildWhere()
	var cols []string
	var args []interface{}
	m := p.model.Obj2Map(obj)
	for k, v := range m {
		cols = append(cols, p.quote(k)+"=?")
		args = append(args, v)
	}
	args = append(args, wargs...)

	sql := fmt.Sprintf("UPDATE %s SET %v %v", p.quote(p.model.Name), strings.Join(cols, ","), wsql)
	return sql, args
}

func (p *Scope) quote(str string) string { return p.orm.dialect.Quote(str) }
func (p *Scope) checkModel(model T) {
	if p.model == nil {
		m := getModelInfo(model)
		p.model = &m
	}
}
func (p *Scope) setWhereIdIfNoWhere(model T) {
	p.checkModel(model)
	if !p.haswhere {
		p.WhereId(p.model.IdVal(model))
	}
}

func (p *Scope) exec(sql string, args ...interface{}) {
	params:=convertArgs(args...)
	if _, err := p.orm.db.Exec(sql, params...); err != nil {
		fmt.Println(sql, args)
		fmt.Println(err)
	}
}

func convertArgs(args... interface{}) []interface{}{
	params := []interface{}{}
	for _, arg := range args {
		switch a := arg.(type) {
		case GUID:
			params = append(params, string(a))
		default:
			params = append(params, a)
		}
	}
	return params
}