package db

import (
	. "github.com/eynstudio/gobreak"
)

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

type PageFilter struct {
	FilterGroup
	Page    int
	PerPage int
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
