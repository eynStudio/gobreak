package db

import (
	. "github.com/eynstudio/gobreak"
)

type SqlArgs struct {
	Sql  string
	Args []interface{}
}

func (p *SqlArgs) AddArgs(a ...interface{}) { p.Args = append(p.Args, a...) }

type Paging struct {
	Total int
	Items T
}
