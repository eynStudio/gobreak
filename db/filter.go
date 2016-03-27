package db

import (
	"fmt"
	"strings"

	. "github.com/eynstudio/gobreak"
)

type VisitorResult struct {
	Sql  string
	Args []interface{}
}
type VisitorResults []VisitorResult

func (p VisitorResults) Join(con string) (vr VisitorResult) {
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

func (p *FilterVisitor) Visitor(filter FilterGroup) (wsql string, args []interface{}) {
	r := p.VisitGroup(filter)
	return r.Sql, r.Args
}

func (p *FilterVisitor) VisitGroup(f FilterGroup) (vr VisitorResult) {
	var lst VisitorResults

	ruleResult := p.VisitRules(f.Con, f.Rules)
	lst = append(lst, ruleResult)

	for _, it := range f.Groups {
		lst = append(lst, p.VisitGroup(it))
	}

	return lst.Join(f.Con)
}

func (p *FilterVisitor) VisitRules(con string, rules []FilterRule) (vr VisitorResult) {
	var lst VisitorResults
	for _, it := range rules {
		lst = append(lst, p.VisitRule(it))
	}
	return lst.Join(con)
}

func (p *FilterVisitor) VisitRule(f FilterRule) (vr VisitorResult) {
	if f.O == "like" {
		vr.Sql = fmt.Sprintf("%s %s ?", f.F, f.O)
		vr.Args = append(vr.Args, "%"+f.V1.(string)+"%")
	} else if f.O == "=" {
		vr.Sql = fmt.Sprintf("%s %s ?", f.F, f.O)
		vr.Args = append(vr.Args, f.V1)
	} else if f.O == "in" {
		var ss []string
		for _, it := range f.V1.([]string) {
			ss = append(ss, "?")
			vr.Args = append(vr.Args, it)
		}
		vs := strings.Join(ss, ",")
		vr.Sql = fmt.Sprintf("%s %s (%s)", f.F, f.O, vs)
	}
	return
}

type Filter interface {
}

type FilterRule struct {
	F  string
	O  string
	V1 T
	V2 T
}

type FilterGroup struct {
	Con    string
	Rules  []FilterRule
	Groups []FilterGroup
}

func (p *FilterGroup) AddRule(r FilterRule) {
	p.Rules = append(p.Rules, r)
}
func (p *FilterGroup) AddGroup(g FilterGroup) {
	p.Groups = append(p.Groups, g)
}

func NewAndFilterGroup() (fg FilterGroup) {
	return FilterGroup{Con: "and"}
}
func NewOrFilterGroup() (fg FilterGroup) {
	return FilterGroup{Con: "or"}
}

type PageFilter struct {
	FilterGroup
	Page    int
	PerPage int
	Search  string
}

func NewPageFilter(page, perPage int, field string, val1 T) *PageFilter {
	p := &PageFilter{
		Page:    page,
		PerPage: perPage,
	}
	p.Rules = append(p.Rules, FilterRule{
		F:  field,
		V1: val1,
	})
	return p
}

func (p *PageFilter) Skip() int {
	return (p.Page - 1) * p.PerPage
}

type Paging struct {
	Total int
	Items T
}
