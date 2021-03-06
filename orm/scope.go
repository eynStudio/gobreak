package orm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db"
	"github.com/eynstudio/gobreak/db/filter"
)

type Scope struct {
	_select string
	_from   string
	Error
	orm       *Orm
	model     *model
	haswhere  bool
	where     map[string]interface{}
	whereid   interface{}
	wheresql  string
	whereargs []interface{}
	pf        *filter.PageFilter
	offset    int
	limit     int
	hasLimit  bool
	*sql.Tx
	builder Ibuilder
}

func NewScope(orm *Orm) *Scope {
	s := &Scope{orm: orm}
	s.builder = orm.getBuilder(s)
	return s
}

func (p *Scope) getSelect() string { return "select " + IfThenStr(p._select == "", "*", p._select) }
func (p *Scope) getFrom() string   { return " from " + p.getTblName() }
func (p *Scope) getTblName() string {
	if p._from != "" {
		return p.quote(p._from)
	}
	return p.quote(p.model.Name)
	//		return p.quote(IfThenStr(p._from == "", p.model.Name, p._from))
}
func (p *Scope) Select(s ...string) *Scope {
	p.builder.Select(s...)
	p._select = strings.Join(s, ",")
	return p
}

func (p *Scope) From(name string) *Scope {
	p.builder.From(name)
	p._from = name
	return p
}
func (p *Scope) Order(args ...string) *Scope {
	p.builder.Order(args...)
	return p
}

func (p *Scope) Where(sql string, args ...interface{}) *Scope {
	p.builder.Where(sql, args...)
	p.haswhere = true
	p.wheresql = sql
	p.whereargs = args
	return p
}

func (p *Scope) WhereId(id interface{}) *Scope {
	p.builder.WhereId(id)
	p.haswhere = true
	p.whereid = id
	return p
}

func (p *Scope) Count(model T) int {
	p.checkModel(model)
	n := ""
	if p.model != nil {
		n = p.model.Name
	}
	sa := p.builder.From(n).SqlCount()
	row := p._queryRow2(sa)
	count := 0
	if err := row.Scan(&count); err != nil {
		log.Println(err)
	}
	return count
}

func (p *Scope) Has(model T) bool {
	return p.Count(model) > 0
}

func (p *Scope) Get(model T) bool {
	p.checkModel(model)
	n := ""
	if p.model != nil {
		n = p.model.Name
	}
	sa := p.builder.From(n).SqlSelect()
	var rows *sql.Rows
	if rows, p.Err = p._query2(sa); p.IsErr() {
		log.Println(p.Err)
		return false
	}
	defer rows.Close()
	if rows.Next() {
		p.model.MapRowsAsObj(rows, model)
		return true
	}
	return false
}

//func (p *Scope) Get3(model T) bool {
//	p.checkModel(model)
//	sa := p.buildQuery()
//	var rows *sql.Rows
//	if rows, p.Err = p._query(sa.Sql, convertArgs(sa)...); p.IsErr() {
//		return false
//	}
//	defer rows.Close()
//	if rows.Next() {
//		p.model.MapRowsAsObj(rows, model)
//		return true
//	}
//	return false
//}

func (p *Scope) buildQuery() (sa db.SqlArgs) {
	sa.Sql += p.getSelect() + p.getFrom()
	sa2 := p.buildWhere()
	sa.Sql += sa2.Sql
	sa.Args = sa2.Args
	return
}

func (p *Scope) All(model T) *Scope {
	p.checkModel(model)
	n := ""
	if p.model != nil {
		n = p.model.Name
	}
	sa := p.builder.From(n).SqlSelect()
	var rows *sql.Rows
	if rows, p.Err = p._query2(sa); p.NotErr() {
		defer rows.Close()
		p.model.MapRowsAsLst(rows, model)
	}
	return p
}

func (p *Scope) Int() (n int) {
	sa := p.builder.SqlSelect()
	r := p._queryRow2(sa)
	if err := r.Scan(&n); err != nil {
		log.Println(err)
	}
	return
}

func (p *Scope) Strs() (lst []string) {
	sa := p.builder.SqlSelect()
	var rows *sql.Rows
	if rows, p.Err = p._query2(sa); p.NotErr() {
		defer rows.Close()
		for rows.Next() {
			var v string
			rows.Scan(&v)
			lst = append(lst, v)
		}
	}
	return
}

func (p *Scope) Guids() (lst []GUID) {
	sa := p.builder.SqlSelect()
	var rows *sql.Rows
	if rows, p.Err = p._query2(sa); p.NotErr() {
		defer rows.Close()
		for rows.Next() {
			var v string
			rows.Scan(&v)
			lst = append(lst, GUID(v))
		}
	}
	return
}

func (p *Scope) Query(model T, query string, args ...interface{}) *Scope {
	p.checkModel(model)
	var rows *sql.Rows

	if rows, p.Err = p._query(query, convertArgs2(args)...); p.NotErr() {
		defer rows.Close()
		p.model.MapRowsAsLst(rows, model)
	}
	return p
}

func (p *Scope) Limit(offset, limit int) *Scope {
	p.builder.Limit(limit, offset)
	p.offset = offset
	p.limit = limit
	p.hasLimit = true
	return p
}

func (p *Scope) Page2(model T, pf *filter.PageFilter) *db.Paging {
	p.checkModel(model)

	if pf != nil {
		p.builder.Limit(pf.PerPage(), pf.Skip())
		visitor := filter.Visitor{}
		visitor.Quote = p.orm.dialect.Quote
		sa := visitor.Visitor(pf.Group)
		p.builder.Where(sa.Sql, sa.Args...)
	}

	var rows *sql.Rows
	paging := &db.Paging{}
	paging.Total = p.Count(model)
	sa := p.builder.SqlSelect()
	if rows, p.Err = p._query2(sa); p.NotErr() {
		defer rows.Close()
		p.model.MapRowsAsLst(rows, model)
		paging.Items = model
	}
	return paging
}
func (p *Scope) PageJson2(lst T, pf *filter.PageFilter) *db.Paging {
	p.checkModel(lst)

	if pf != nil {
		p.builder.Limit(pf.PerPage(), pf.Skip())
		visitor := filter.Visitor{}
		visitor.Quote = p.orm.dialect.Quote
		sa := visitor.Visitor(pf.Group)
		p.builder.Where(sa.Sql, sa.Args...)
	}

	var rows *sql.Rows
	paging := &db.Paging{}
	paging.Total = p.Count(lst)

	if !p.builder.hasSelect() {
		p.builder.Select(p.orm.Mapper().Map2Db("Json"))
	}
	sa := p.builder.SqlSelect()

	resultv := reflect.ValueOf(lst)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("out argument must be a slice address")
	}
	slicev := resultv.Elem()
	if rows, p.Err = p._query2(sa); p.NotErr() {
		defer rows.Close()
		for rows.Next() {
			var v []byte
			rows.Scan(&v)
			obj := reflect.New(p.model.Type).Interface()
			json.Unmarshal(v, obj)
			slicev = reflect.Append(slicev, reflect.ValueOf(obj).Elem())
		}
		resultv.Elem().Set(slicev.Slice(0, slicev.Len()))
	}
	paging.Items = lst
	return paging
}
func (p *Scope) Page(model T, pf *filter.PageFilter) *db.Paging {
	p.checkModel(model)
	p.pf = pf
	p.haswhere = true
	p.Limit(pf.Skip(), pf.PerPage())
	w := p.buildWhere()
	psa := p.buildPage()
	sql_ := fmt.Sprintf("SELECT * %s %v %v", p.getFrom(), w.Sql, psa.Sql)
	var rows *sql.Rows
	log.Println(sql_, w.Args)
	paging := &db.Paging{}
	if rows, p.Err = p._query(sql_, convertArgs(w)...); p.NotErr() {
		defer rows.Close()

		p.model.MapRowsAsLst(rows, model)
		paging.Total = p.Count(model)
		paging.Items = model
	}
	return paging
}

func (p *Scope) PageByOrder(model T, order string, pf *filter.PageFilter) *db.Paging {
	p.checkModel(model)
	p.pf = pf
	p.haswhere = true
	p.Limit(pf.Skip(), pf.PerPage())
	w := p.buildWhere()
	psa := p.buildPageByOrder(order)
	sql_ := fmt.Sprintf("SELECT * %s %v %v", p.getFrom(), w.Sql, psa.Sql)
	var rows *sql.Rows
	paging := &db.Paging{}
	log.Println(sql_)
	if rows, p.Err = p._query(sql_, convertArgs(w)...); p.NotErr() {
		defer rows.Close()

		p.model.MapRowsAsLst(rows, model)
		paging.Total = p.Count(model)
		paging.Items = model
	}
	return paging
}

//func (p *Scope) PageSql(model T, pf filter.PageFilter, sql string) *db.Paging {
//	p.checkModel(model)
//	paging := &db.Paging{}
//	//sql="select "

//	return paging
//}
func (p *Scope) Save(model T) *Scope {
	p.checkModel(model)
	p.setWhereIdIfNoWhere(model)

	if p.orm.dialect.Driver() == "postgres" {
		return p.exec(p.buildUpsert(model))
	}

	if p.Has(model) {
		p.Update(model)
	} else {
		p.Insert(model)
	}
	return p
}

func (p *Scope) SaveTo(name string, model T) *Scope {
	p._from = name
	p.checkModel(model)

	if p.orm.dialect.Driver() == "postgres" {
		return p.exec(p.buildUpsert(model))
	}

	if p.Has(model) {
		p.Update(model)
	} else {
		p.Insert(model)
	}
	return p
}

func (p *Scope) SaveJson(id GUID, data T) *Scope {
	p.checkModel(data)
	sa := p.builder.From(p.model.Name).SqlSaveJson(id, data)
	p.exec(*sa)
	return p
}

func (p *Scope) GetJson2(data T) bool {
	p.checkModel(data)
	p.Select(p.orm.Mapper().Map2Db("Json")).From(p.model.Name)
	sa := p.builder.SqlSelect()
	rows, err := p._query2(sa)
	if p.Err = err; p.IsErr() {
		log.Println(p.Err)
		return false
	}
	defer rows.Close()

	if !rows.Next() {
		return false
	}

	var vv []byte
	p.Err = rows.Scan(&vv)
	p.NoErrExec(func() { p.Err = json.Unmarshal(vv, &data) })
	return !p.IsErr()
}

func (p *Scope) GetJson(data T) bool {
	p.checkModel(data)
	p.Select(`"Json"`)
	sa := p.buildQuery()
	rows, err := p._query(sa.Sql, convertArgs(sa)...)
	if p.Err = err; p.IsErr() {
		log.Println(p.Err)
		return false
	}
	defer rows.Close()

	if !rows.Next() {
		return false
	}

	var vv []byte
	p.Err = rows.Scan(&vv)
	p.NoErrExec(func() { p.Err = json.Unmarshal(vv, &data) })
	return true
}

func (p *Scope) AllJson(lst T) *Scope {
	p.checkModel(lst)
	p.Select(p.orm.Mapper().Map2Db("Json"))
	p.From(p.model.Name)
	resultv := reflect.ValueOf(lst)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("out argument must be a slice address")
	}
	slicev := resultv.Elem()

	sa := p.builder.SqlSelect()

	var rows *sql.Rows
	if rows, p.Err = p._query2(sa); p.NotErr() {
		defer rows.Close()
		for rows.Next() {
			var v []byte
			rows.Scan(&v)
			obj := reflect.New(p.model.Type).Interface()
			json.Unmarshal(v, obj)
			slicev = reflect.Append(slicev, reflect.ValueOf(obj).Elem())
		}
		resultv.Elem().Set(slicev.Slice(0, slicev.Len()))
	}
	return p
}

//func (p *Scope) AllJson(lst T) *Scope {
//	p.checkModel(lst)
//	resultv := reflect.ValueOf(lst)
//	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
//		panic("out argument must be a slice address")
//	}
//	slicev := resultv.Elem()

//	sql_ := fmt.Sprintf(`SELECT "Json" %s where `+p.wheresql, p.getFrom())
//	var rows *sql.Rows
//	if rows, p.Err = p._query(sql_, p.whereargs...); p.NotErr() {
//		defer rows.Close()
//		for rows.Next() {
//			var v []byte
//			rows.Scan(&v)
//			obj := reflect.New(p.model.Type).Interface()
//			json.Unmarshal(v, obj)
//			slicev = reflect.Append(slicev, reflect.ValueOf(obj).Elem())
//		}
//		resultv.Elem().Set(slicev.Slice(0, slicev.Len()))
//	}
//	return p
//}
func (p *Scope) PageJson(lst T, page, perPage int) (pager db.Paging) {
	p.checkModel(lst)
	pf := filter.NewPageFilter(page, perPage)
	p.Limit(pf.Skip(), perPage)
	pager.Total = p.Count(lst)
	//	log.Println(pager.Total)
	resultv := reflect.ValueOf(lst)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("out argument must be a slice address")
	}
	slicev := resultv.Elem()
	wsa := p.buildWhere()
	pas := p.buildPage()
	sql_ := fmt.Sprintf(`SELECT "Json" %s `+wsa.Sql+" "+pas.Sql, p.getFrom())
	var rows *sql.Rows
	if rows, p.Err = p._query(sql_, wsa.Args...); p.NotErr() {
		defer rows.Close()
		for rows.Next() {
			var v []byte
			rows.Scan(&v)
			obj := reflect.New(p.model.Type).Interface()
			json.Unmarshal(v, obj)
			slicev = reflect.Append(slicev, reflect.ValueOf(obj).Elem())
		}
		resultv.Elem().Set(slicev.Slice(0, slicev.Len()))
	}
	pager.Items = lst
	//	log.Println(lst)
	return
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
	sa := p.buildUpdateFields(model, fields)
	p.exec(sa)
	return p
}

func (p *Scope) Del(model T) *Scope {
	p.checkModel(model)
	p.From(p.model.Name)
	p.setWhereIdIfNoWhere(model)
	sa := p.builder.SqlDel()

	return p.exec(*sa)
	//	w := p.buildWhere()
	//	w.Sql = fmt.Sprintf("DELETE from %s %v", p.quote(p.model.Name), w.Sql)
	//	p.exec(w)
	//	return p
}

func (p *Scope) DelAll(model T) *Scope {
	p.checkModel(model)
	p.From(p.model.Name)
	sa := p.builder.SqlDel()
	return p.exec(*sa)
	//	var sa db.SqlArgs
	//	sa.Sql = fmt.Sprintf("DELETE from %s", p.quote(p.model.Name))
	//	p.exec(sa)
	//	return p
}

func (p *Scope) buildWhere() (sa db.SqlArgs) {
	if !p.haswhere {
		return
	} else if p.whereid != nil {
		sa.Sql = fmt.Sprintf(" WHERE (%v=?)", p.quote(p.model.Id()))
		sa.AddArgs(p.whereid)
	} else if len(p.wheresql) > 0 {
		sa.Sql = " WHERE " + p.wheresql
		sa.AddArgs(p.whereargs...)
	} else if p.pf != nil {
		visitor := filter.Visitor{}
		visitor.Quote = p.orm.dialect.Quote
		sa = visitor.Visitor(p.pf.Group)
		if sa.Sql != "" {
			sa.Sql = " WHERE " + sa.Sql
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
	if p.orm.dialect.Driver() == "postgres" {
		sa.Sql = fmt.Sprintf("limit %v offset %v", p.limit, p.offset)
		//		log.Println(sa.Sql)
	} else if p.orm.dialect.Driver() == "mysql" {
		sa.Sql = fmt.Sprintf("limit %v,%v", p.offset, p.limit)
	} else if p.orm.dialect.Driver() == "mssql" {
		sa.Sql = fmt.Sprintf("ORDER BY %v OFFSET %v ROW FETCH NEXT %v ROWS only", p.model.Id(), p.offset, p.limit)
	}

	return
}

func (p *Scope) buildPageByOrder(order string) (sa db.SqlArgs) {
	if !p.hasLimit {
		return
	}

	if p.orm.dialect.Driver() == "postgres" {
		sa.Sql = fmt.Sprintf("order by %v limit %v offset %v", order, p.limit, p.offset)
		log.Println(sa.Sql)
	} else if p.orm.dialect.Driver() == "mysql" {
		sa.Sql = fmt.Sprintf("order by %v limit %v,%v", order, p.offset, p.limit)
	} else if p.orm.dialect.Driver() == "mssql" {
		sa.Sql = fmt.Sprintf("ORDER BY %v OFFSET %v ROW FETCH NEXT %v ROWS only", order, p.offset, p.limit)
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
	sa.Sql = fmt.Sprintf("insert into %s (%v) values (%v)", p.getTblName(), strings.Join(cols, ","), strings.Join(params, ","))
	return
}

func (p *Scope) buildUpsert(obj T) (sa db.SqlArgs) {
	var cols []string
	var params []string
	m := p.model.Obj2Map(obj)
	i := 1
	for k, v := range m {
		cols = append(cols, p.quote(k))
		params = append(params, fmt.Sprintf("$%d", i))
		sa.AddArgs(v)
		i += 1
	}
	cc := strings.Join(cols, ",")
	pp := strings.Join(params, ",")
	sa.Sql = fmt.Sprintf(`insert into %s (%v) values (%v) ON CONFLICT ("Id") DO UPDATE SET (%v)=(%v)`,
		p.getTblName(), cc, pp, cc, pp,
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

	sa.Sql = fmt.Sprintf("UPDATE %s SET %v %v", p.getTblName(), strings.Join(cols, ","), w.Sql)
	return
}

func (p *Scope) buildUpdateFields(obj T, fields []string) (sa db.SqlArgs) {
	w := p.buildWhere()
	var cols []string
	m := p.model.Obj2Map(obj)
	for k, v := range m {
		for _, it := range fields {
			if it == k {
				cols = append(cols, p.quote(k)+"=?")
				sa.AddArgs(v)
			}
		}
	}
	sa.AddArgs(w.Args...)

	sa.Sql = fmt.Sprintf("UPDATE %s SET %v %v", p.getTblName(), strings.Join(cols, ","), w.Sql)
	return
}

func (p *Scope) quote(str string) string { return p.orm.dialect.Quote(str) }
func (p *Scope) checkModel(model T) {
	if p.model == nil && model != nil {
		p.model = p.orm.models.get(model)
	}
}
func (p *Scope) setWhereIdIfNoWhere(model T) {
	p.checkModel(model)
	if !p.haswhere {
		p.WhereId(p.model.IdVal(model))
	}
}

func convertArgs(sa db.SqlArgs) []interface{} {
	return convertArgs2(sa.Args)
}

func convertArgs2(args []interface{}) []interface{} {
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

func (p Scope) IsNotFound() bool { return p.IsErr() && p.Err == db.DbNotFound }

func (p *Scope) _query2(sa *db.SqlArgs) (*sql.Rows, error) {
	query := p.orm.convParams(sa.Sql)
	args := convertArgs2(sa.Args)
	//	log.Println(query, args)
	if p.hasTx() {
		return p.Tx.Query(query, args...)
	}
	return p.orm.db.Query(query, args...)
}
func (p *Scope) _queryRow2(sa *db.SqlArgs) *sql.Row {
	query := p.orm.convParams(sa.Sql)
	args := convertArgs2(sa.Args)
	//	log.Println(query, args)

	if p.hasTx() {
		return p.Tx.QueryRow(query, args...)
	}
	return p.orm.db.QueryRow(query, args...)
}

func (p *Scope) _query(query string, args ...interface{}) (*sql.Rows, error) {
	query = p.orm.convParams(query)
	if p.hasTx() {
		return p.Tx.Query(query, args...)
	}
	return p.orm.db.Query(query, args...)
}

func (p *Scope) _queryRow(query string, args ...interface{}) *sql.Row {
	query = p.orm.convParams(query)
	//	log.Println(query, args)

	if p.hasTx() {
		return p.Tx.QueryRow(query, args...)
	}
	return p.orm.db.QueryRow(query, args...)
}

func (p *Scope) exec(sa db.SqlArgs) *Scope {
	params := convertArgs(sa)
	query := p.orm.convParams(sa.Sql)
	//	log.Println(query, params)

	if p.hasTx() {
		_, p.Err = p.Tx.Exec(query, params...)

	} else {
		_, p.Err = p.orm.db.Exec(query, params...)
	}
	return p
}

func (p Scope) hasTx() bool { return p.Tx != nil }
