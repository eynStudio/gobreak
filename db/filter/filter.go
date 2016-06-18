package filter

import (
	. "github.com/eynstudio/gobreak"
)

type Rule struct {
	F  string
	O  string
	V1 string
	V2 string
}

type Group struct {
	Con    string
	Rules  []Rule
	Groups []Group
}

func (p *Group) AddRule(r Rule)   { p.Rules = append(p.Rules, r) }
func (p *Group) AddGroup(g Group) { p.Groups = append(p.Groups, g) }
func NewAndGroup() (fg Group)     { return Group{Con: "and"} }
func NewOrGroup() (fg Group)      { return Group{Con: "or"} }

type Filter struct {
	Group
	Ext M
}

type PageFilter struct{ Filter }

func (p *Filter) Search() string { return p.Ext.GetStr("search") }
func (p *Filter) Role() string   { return p.Ext.GetStr("role") }
func (p *PageFilter) Page() int {
	n := p.Ext.GetIntOr("page", 1)
	if n < 1 {
		n = 1
	}
	return n
}
func (p *PageFilter) PerPage() int {
	n := p.Ext.GetIntOr("perPage", 20)
	if n < 1 {
		n = 20
	}
	return n
}
func (p *PageFilter) Skip() int { return (p.Page() - 1) * p.PerPage() }

func NewPageFilter(page, perPage int) *PageFilter {
	p := &PageFilter{}
	p.Ext = M{"page": page, "perPage": perPage}
	return p
}
