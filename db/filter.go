package db

import (
	"fmt"
	"strings"

	. "github.com/eynstudio/gobreak"
)

type FilterVisitor struct {
}

func (p *FilterVisitor) Visitor(filter FilterGroup) (wsql string, args []interface{}) {
	return p.VisitGroup(filter)
}

func (p *FilterVisitor) VisitGroup(f FilterGroup) (wsql string, args []interface{}) {
	if f.Con == "and" || f.Con == "or" {
		if len(f.Rules) == 0 {
			return
		} else if len(f.Rules) == 1 {
			return p.VisitRule(f.Rules[0])
		} else {
			wsql, args = p.VisitRule(f.Rules[0])
			for _, it := range f.Rules[1:] {
				w, a := p.VisitRule(it)
				wsql = wsql + " " + f.Con + " " + w
				args = append(args, a...)
			}
			return
		}

	} else if f.Con == "or" {

	}

	return
}

func (p *FilterVisitor) VisitRule(f FilterRule) (wsql string, args []interface{}) {
	if f.O == "like" {
		return fmt.Sprintf("%s %s ?", f.F, f.O), append(args, "%"+f.V1.(string)+"%")
	} else if f.O == "=" {
		return fmt.Sprintf("%s %s ?", f.F, f.O), append(args, f.V1)
	} else if f.O == "in" {
		var ss []string
		for _, it := range f.V1.([]string) {
			ss = append(ss, "?")
			args = append(args, it)
		}
		vs := strings.Join(ss, ",")
		return fmt.Sprintf("%s %s (%s)", f.F, f.O, vs), args
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
