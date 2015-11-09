package orm

import (
	"fmt"

	. "github.com/eynstudio/gobreak"
)

type Scope struct {
	orm       *Orm
	model     model
	where     map[string]interface{}
	whereid   interface{}
	wheresql  string
	whereargs []interface{}
}

func NewScope(orm *Orm) *Scope {
	return &Scope{orm: orm}
}

func (p *Scope) Where(sql string, args ...interface{}) *Scope {
	p.wheresql = sql
	p.whereargs = args
	return p
}

func (p *Scope) WhereId(id interface{}) *Scope {
	p.whereid = id
	return p
}

func (p *Scope) Count(model T) int {
	p.model = p.orm.models.GetModelInfo(model)
	id := p.model.Id()
	wsql, args := p.buildWhere()
	sql := fmt.Sprintf("SELECT COUNT(%v) from %s %v", p.quote(id), p.quote(p.model.Name), wsql)
	row := p.orm.db.QueryRow(sql, args...)
	count := 0
	row.Scan(&count)
	return count
}

func (p *Scope) One(model T) (has bool) {
	p.model = p.orm.models.GetModelInfo(model)
	wsql, args := p.buildWhere()
	sql := fmt.Sprintf("SELECT TOP 1 * from %s %v", p.quote(p.model.Name), wsql)
	rows, _ := p.orm.db.Query(sql, args...)
	p.model.MapRowsAsObj(rows, model)
	return true
}

func (p *Scope) All(model T) {
	p.model = p.orm.models.GetModelInfo(model)
	wsql, args := p.buildWhere()
	sql := fmt.Sprintf("SELECT * from %s %v", p.quote(p.model.Name), wsql)
	rows, _ := p.orm.db.Query(sql, args...)
	p.model.MapRowsAsLst(rows, model)
}

func (p *Scope) buildWhere() (string, []interface{}) {
	if p.whereid != nil {
		return fmt.Sprintf("WHERE (%v=?)", p.model.Id()), []interface{}{p.whereid}
	} else if len(p.wheresql) > 0 {
		return "WHERE " + p.wheresql, p.whereargs
	} else if len(p.where) > 0 {

	}

	return "", nil
}
func (p *Scope) quote(str string) string { return p.orm.dialect.Quote(str) }
