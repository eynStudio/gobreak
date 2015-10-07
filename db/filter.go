package db

import (
	. "github.com/eynstudio/gobreak"
)

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

type PageFilter struct {
	FilterGroup
	Page    int
	PerPage int
}

func (p *PageFilter) Skip() int{
	return (p.Page-1)*p.PerPage
}
type Paging struct {
	Total   int
	Items   T
}
