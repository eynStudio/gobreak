package db

import (
	"errors"

	. "github.com/eynstudio/gobreak"
)

var DbNotFound = errors.New("Db Not Found")

type SqlArgs struct {
	Sql  string
	Args []interface{}
}

func NewAgrs(sql string, args ...interface{}) *SqlArgs { return &SqlArgs{Sql: sql, Args: args} }

func (p *SqlArgs) AddArgs(a ...interface{}) { p.Args = append(p.Args, a...) }

func (p *SqlArgs) Append(sql string, args ...interface{}) *SqlArgs {
	return NewAgrs(p.Sql+sql, append(p.Args, args...)...)
}

func (p *SqlArgs) Append2(sa *SqlArgs) *SqlArgs {
	if sa == nil {
		return NewAgrs(p.Sql, p.Args...)
	}
	return NewAgrs(p.Sql+sa.Sql, append(p.Args, sa.Args...)...)
}

type Paging struct {
	Total int
	Items T
}
