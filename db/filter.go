package db

import (
	"fmt"
	"strings"

	. "github.com/eynstudio/gobreak"
)

type VisitorResults []SqlArgs

func (p VisitorResults) Join(con string) (vr SqlArgs) {
	if con == "" {
		con = "and"
	}

	var sqls []string
	for _, it := range p {
		if it.Sql != "" {
			sqls = append(sqls, it.Sql)
			vr.Args = append(vr.Args, it.Args...)
		}
	}
	if len(sqls) > 0 {
		vr.Sql = "(" + strings.Join(sqls, " "+con+" ") + ")"
	}
	return
}

type FilterVisitor struct {
}

func (p *FilterVisitor) Visitor(filter FilterGroup) SqlArgs {
	return p.VisitGroup(filter)
}

func (p *FilterVisitor) VisitGroup(f FilterGroup) (vr SqlArgs) {
	var lst VisitorResults

	ruleResult := p.VisitRules(f.Con, f.Rules)
	lst = append(lst, ruleResult)

	for _, it := range f.Groups {
		lst = append(lst, p.VisitGroup(it))
	}

	return lst.Join(f.Con)
}

func (p *FilterVisitor) VisitRules(con string, rules []FilterRule) (vr SqlArgs) {
	var lst VisitorResults
	for _, it := range rules {
		lst = append(lst, p.VisitRule(it))
	}
	return lst.Join(con)
}

func (p *FilterVisitor) VisitRule(f FilterRule) (vr SqlArgs) {
	if f.O == "like" {
		vr.Sql = fmt.Sprintf("%s %s ?", f.F, f.O)
		vr.Args = append(vr.Args, "%"+f.V1+"%")
	} else if f.O == "=" {
		vr.Sql = fmt.Sprintf("%s %s ?", f.F, f.O)
		vr.Args = append(vr.Args, f.V1)
	} else if f.O == "in" {
		var ss []string
		vlst := strings.Split(f.V1, ",")
		for _, it := range vlst {
			ss = append(ss, "?")
			vr.Args = append(vr.Args, it)
		}
		vs := strings.Join(ss, ",")
		vr.Sql = fmt.Sprintf("%s %s (%s)", f.F, f.O, vs)
	}
	return
}

type FilterRule struct {
	F  string
	O  string
	V1 string
	V2 string
}

type FilterGroup struct {
	Con    string
	Rules  []FilterRule
	Groups []FilterGroup
}

func (p *FilterGroup) AddRule(r FilterRule)   { p.Rules = append(p.Rules, r) }
func (p *FilterGroup) AddGroup(g FilterGroup) { p.Groups = append(p.Groups, g) }
func NewAndFilterGroup() (fg FilterGroup)     { return FilterGroup{Con: "and"} }
func NewOrFilterGroup() (fg FilterGroup)      { return FilterGroup{Con: "or"} }

type Filter struct {
	FilterGroup
	Ext M
}

type PageFilter struct{ Filter }

func (p *Filter) Search() string   { return p.Ext.GetStr("search") }
func (p *Filter) Role() string     { return p.Ext.GetStr("role") }
func (p *PageFilter) Page() int    { return int(p.Ext.GetInt("page")) }
func (p *PageFilter) PerPage() int { return int(p.Ext.GetInt("perPage")) }
func (p *PageFilter) Skip() int    { return (p.Page() - 1) * p.PerPage() }

func NewPageFilter(page, perPage int, field string, val1 string) *PageFilter {
	p := &PageFilter{}
	p.Ext = M{"page": page, "perPage": perPage}
	p.Rules = append(p.Rules, FilterRule{F: field, V1: val1})
	return p
}

type Paging struct {
	Total int
	Items T
}
