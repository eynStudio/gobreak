package orm

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db"
)

type Ibuilder interface {
	Where(sql string, args ...interface{}) Ibuilder
	WhereId(id interface{}) Ibuilder
	Order(args ...string) Ibuilder
	Limit(n, offset int) Ibuilder
	Select(f ...string) Ibuilder
	From(f string) Ibuilder
	SqlSelect() (sa *db.SqlArgs)
	SqlDel() (sa *db.SqlArgs)
	SqlCount() (sa *db.SqlArgs)
	SqlSaveJson(id GUID, data T) (sa *db.SqlArgs)
	hasSelect() bool
	hasWhere() bool
}

type builder struct {
	limitArgs *db.SqlArgs
	whereArgs *db.SqlArgs
	orders    []string
	fields    []string
	from      string
	mapper    MapperFn
	scope     *Scope
}

func newBuilder(s *Scope) *builder {
	b := &builder{scope: s}
	b.mapper = s.orm.mapper
	return b
}
func (p builder) hasLimit() bool  { return p.limitArgs != nil }
func (p builder) hasWhere() bool  { return p.whereArgs != nil }
func (p builder) hasOrder() bool  { return len(p.orders) > 0 }
func (p builder) hasSelect() bool { return len(p.fields) > 0 }

func (p *builder) From(f string) Ibuilder {
	if p.from == "" {
		p.from = p.mapper(f)
	}
	return p
}
func (p *builder) Select(f ...string) Ibuilder {
	p.fields = append(p.fields, f...)
	return p
}
func (p *builder) Where(sql string, args ...interface{}) Ibuilder {
	if sql == "" {
		return p
	}
	p.whereArgs = p.initWhereArgs().Append(sql, args...)
	return p
}
func (p *builder) WhereId(id interface{}) Ibuilder {
	p.whereArgs = p.initWhereArgs().Append(p.mapper("Id")+"=?", id)
	return p
}
func (p *builder) initWhereArgs() *db.SqlArgs {
	if p.whereArgs == nil {
		p.whereArgs = db.NewAgrs(" WHERE ")
	}
	return p.whereArgs
}
func (p *builder) Order(args ...string) Ibuilder {
	p.orders = append(p.orders, args...)
	return p
}
func (p *builder) Limit(n, offset int) Ibuilder {
	p.limitArgs = db.NewAgrs(` LIMIT ? OFFSET ?`, n, offset)
	return p
}

func (p *builder) SqlSelect() (sa *db.SqlArgs) {
	sql := `SELECT ` + p.buildFields() + " FROM " + p.mapper(p.from)
	sa = db.NewAgrs(sql)
	sa = sa.Append2(p.whereArgs)
	if p.hasOrder() {
		orders := " ORDER BY " + strings.Join(p.orders, ",")
		sa = sa.Append(orders)
	}
	return sa.Append2(p.limitArgs)
}

func (p *builder) SqlDel() (sa *db.SqlArgs) {
	sql := `DELETE FROM ` + p.mapper(p.from)
	sa = db.NewAgrs(sql)
	return sa.Append2(p.whereArgs)
}

func (p *builder) SqlCount() (sa *db.SqlArgs) {
	sql := `SELECT count(` + p.mapper("Id") + ") FROM " + p.mapper(p.from)
	sa = db.NewAgrs(sql)
	return sa.Append2(p.whereArgs)
}

func (p *builder) SqlSaveJson(id GUID, data T) (sa *db.SqlArgs) {
	buf, _ := json.Marshal(data)
	_id, _json := p.mapper("Id"), p.mapper("Json")
	sql := fmt.Sprintf(`Insert into %s(%s,%s) values($1,$2) ON CONFLICT (%s) DO UPDATE SET (%s,%s)=($1,$2)`,
		p.mapper(p.from), _id, _json, _id, _id, _json)

	return db.NewAgrs(sql, id, buf)
}

func (p *builder) buildFields() string {
	if len(p.fields) == 0 {
		return "*"
	}
	return strings.Join(p.fields, ",")
}

type pgBuilder struct {
	*builder
}

type mysqlBuilder struct {
	*builder
}

type oci8Builder struct {
	*builder
}

type mssqlBuilder struct {
	*builder
}

func (p *mssqlBuilder) Limit(n, offset int) Ibuilder {
	p.limitArgs = db.NewAgrs(` OFFSET ? ROW FETCH NEXT ? ROWS only`, offset, n)
	return p
}
