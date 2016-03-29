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
	pf        *db.PageFilter
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
	w := p.buildWhere()
	sql := fmt.Sprintf("SELECT COUNT(%v) from %s %v", p.quote(id), p.quote(p.model.Name), w.Sql)
	row := p.orm.db.QueryRow(sql, convertArgs(w)...)
	count := 0
	row.Scan(&count)
	return count
}
func (p *Scope) Has(model T) bool {
	p.checkModel(model)
	id := p.model.Id()
	var sa db.SqlArgs
	sa.AddArgs(p.model.IdVal(model))
	sa.Sql = fmt.Sprintf("SELECT COUNT(%v) from %s WHERE %v=?", p.quote(id), p.quote(p.model.Name), id)
	row := p.orm.db.QueryRow(sa.Sql, convertArgs(sa)...)
	count := 0
	if err := row.Scan(&count); err != nil {
		fmt.Println("Has:", sa.Sql)
		fmt.Println(err)
	}
	return count > 0
}
func (p *Scope) One(model T) (has bool) {
	p.checkModel(model)
	sa := p.orm.dialect.BulidTopNSql(p, 1)
	rows, err := p.orm.db.Query(sa.Sql, convertArgs(sa)...)

	if err != nil {
		fmt.Println(err, sa.Sql)
	}
	p.model.MapRowsAsObj(rows, model)
	return true
}

func (p *Scope) All(model T) {
	p.checkModel(model)
	w := p.buildWhere()
	sql := fmt.Sprintf("SELECT * from %s %v", p.quote(p.model.Name), w.Sql)
	rows, err := p.orm.db.Query(sql, convertArgs(w)...)
	if err != nil {
		fmt.Println(err, sql)
	}
	p.model.MapRowsAsLst(rows, model)
}

func (p *Scope) Query(model T, query string, args ...interface{}) {
	p.checkModel(model)
	rows, err := p.orm.db.Query(query, args...)
	if err != nil {
		fmt.Println(err, query)
	}
	p.model.MapRowsAsLst(rows, model)
}

func (p *Scope) Limit(offset, limit int) *Scope {
	p.offset = offset
	p.limit = limit
	p.hasLimit = true
	return p
}
func (p *Scope) Page(model T, pf *db.PageFilter) *db.Paging {
	p.checkModel(model)
	p.pf = pf
	p.haswhere = true
	p.Limit(pf.Skip(), pf.PerPage)
	w := p.buildWhere()
	psa := p.buildPage()
	sql := fmt.Sprintf("SELECT * from %s %v %v", p.quote(p.model.Name), w.Sql, psa.Sql)
	rows, err := p.orm.db.Query(sql, convertArgs(w)...)
	if err != nil {
		fmt.Println(sql)
		fmt.Println(convertArgs(w)...)
		panic(err)
	}
	p.model.MapRowsAsLst(rows, model)
	paging := &db.Paging{}
	paging.Total = p.Count(model)
	paging.Items = model
	return paging
}
func (p *Scope) PageSql(model T, pf db.PageFilter, sql string) *db.Paging {
	p.checkModel(model)
	paging := &db.Paging{}
	//sql="select "

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
	sa := p.buildInsert(model)
	p.exec(sa)
	return p
}

func (p *Scope) Update(model T) *Scope {
	p.checkModel(model)
	p.setWhereIdIfNoWhere(model)
	sa := p.buildUpdate(model)
	p.exec(sa)
	return p
}

func (p *Scope) UpdateFields(model T, fields []string) *Scope {
	p.checkModel(model)
	p.setWhereIdIfNoWhere(model)
	sa := p.buildUpdate(model)
	p.exec(sa)
	return p
}

func (p *Scope) Del(model T) *Scope {
	p.checkModel(model)
	p.setWhereIdIfNoWhere(model)
	w := p.buildWhere()
	w.Sql = fmt.Sprintf("DELETE from %s %v", p.quote(p.model.Name), w.Sql)
	p.exec(w)
	return p
}

func (p *Scope) DelAll(model T) *Scope {
	p.checkModel(model)
	var sa db.SqlArgs
	sa.Sql = fmt.Sprintf("DELETE from %s", p.quote(p.model.Name))
	p.exec(sa)
	return p
}

func (p *Scope) buildWhere() (sa db.SqlArgs) {
	if !p.haswhere {
		return
	} else if p.whereid != nil {
		sa.Sql = fmt.Sprintf("WHERE (%v=?)", p.quote(p.model.Id()))
		sa.AddArgs(p.whereid)
	} else if len(p.wheresql) > 0 {
		sa.Sql = "WHERE " + p.wheresql
		sa.AddArgs(p.whereargs...)
	} else if p.pf != nil {
		visitor := db.FilterVisitor{}
		sa = visitor.Visitor(p.pf.FilterGroup)
		if sa.Sql != "" {
			sa.Sql = "WHERE " + sa.Sql
		}
		return
	} else if len(p.where) > 0 {

	}

	return
}

func (p *Scope) buildPage() (sa db.SqlArgs) {
	if !p.hasLimit {
		return
	}

	if p.orm.dialect.Driver() == "mssql" {
		sa.Sql = fmt.Sprintf("ORDER BY %v OFFSET %v ROW FETCH NEXT %v ROWS only", p.model.Id(), p.offset, p.limit)
	}

	return
}
func (p *Scope) buildInsert(obj T) (sa db.SqlArgs) {
	var cols []string
	var params []string
	m := p.model.Obj2Map(obj)
	for k, v := range m {
		cols = append(cols, p.quote(k))
		params = append(params, "?")
		sa.AddArgs(v)
	}
	sa.Sql = fmt.Sprintf("insert into %s (%v) values (%v)",
		p.quote(p.model.Name),
		strings.Join(cols, ","),
		strings.Join(params, ","),
	)
	return
}

func (p *Scope) buildUpdate(obj T) (sa db.SqlArgs) {
	w := p.buildWhere()
	var cols []string
	m := p.model.Obj2Map(obj)
	for k, v := range m {
		cols = append(cols, p.quote(k)+"=?")
		sa.AddArgs(v)
	}
	sa.AddArgs(w.Args...)

	sa.Sql = fmt.Sprintf("UPDATE %s SET %v %v", p.quote(p.model.Name), strings.Join(cols, ","), w.Sql)
	return
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

func (p *Scope) exec(sa db.SqlArgs) {
	params := convertArgs(sa)
	if _, err := p.orm.db.Exec(sa.Sql, params...); err != nil {
		fmt.Println(sa.Sql, sa.Args)
		fmt.Println(err)
	}
}

func convertArgs(sa db.SqlArgs) []interface{} {
	params := []interface{}{}
	for _, arg := range sa.Args {
		switch a := arg.(type) {
		case GUID:
			params = append(params, string(a))
		default:
			params = append(params, a)
		}
	}
	return params
}
