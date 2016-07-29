package orm

import (
	"fmt"

	"github.com/eynstudio/gobreak/db"
)

type Ibuilder interface {
	Where(sql string, args ...interface{}) Ibuilder
	WhereId(id interface{}) Ibuilder
	Order(args ...string) Ibuilder
	Limit(n, offset int) Ibuilder
	Select(f ...string) Ibuilder
	From(f string) Ibuilder
}

type builder struct {
	limit     int
	offset    int
	id        interface{}
	whereArgs *db.SqlArgs
	orders    []string
	fields    []string
	from      string
	mapper    MapperFn
	scope     *Scope
}

func newBuilder(s *Scope) Ibuilder { return &builder{scope: s} }
func (p builder) hasLimit() bool   { return p.limit > 0 }
func (p builder) hasId() bool      { return p.id != nil }
func (p builder) hasOrder() bool   { return len(p.orders) > 0 }

func (p *builder) From(f string) Ibuilder {
	p.from = f
	return p
}
func (p *builder) Select(f ...string) Ibuilder {
	p.fields = append(p.fields, f...)
	return p
}
func (p *builder) Where(sql string, args ...interface{}) Ibuilder {
	p.whereArgs = db.NewAgrs(sql, args...)
	return p
}
func (p *builder) WhereId(id interface{}) Ibuilder {
	p.id = id
	return p
}
func (p *builder) Order(args ...string) Ibuilder {
	p.orders = append(p.orders, args...)
	return p
}
func (p *builder) Limit(n, offset int) Ibuilder {
	p.limit, p.offset = n, offset
	return p
}

func (p *builder) buildWhere() (sa *db.SqlArgs) {
	if p.hasId() {
		return db.NewAgrs(fmt.Sprintf(" WHERE (%v=?)", p.mapper("Id")), p.id)
	}
	return p.whereArgs
}
