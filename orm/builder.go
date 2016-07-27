package orm

import "github.com/eynstudio/gox"

type Ibuilder interface {
	Where(sql string, args ...gox.T) Ibuilder
	WhereId(id gox.T) Ibuilder
	Order(args ...string) Ibuilder
	Limit(n, offset int) Ibuilder
	Select(f ...string) Ibuilder
	From(f string) Ibuilder
}

type builder struct {
	limit  int
	offset int
	id     gox.T
	wsql   string
	args   []gox.T
	orders []string
	fields []string
	from   string
	mapper MapperFn
	scope  *Scope
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
func (p *builder) Where(sql string, args ...gox.T) Ibuilder {
	p.wsql = sql
	p.args = append(p.args, args...)
	return p
}
func (p *builder) WhereId(id gox.T) Ibuilder {
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
