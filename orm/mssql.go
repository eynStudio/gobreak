package orm

import (
	"fmt"
)

type mssql struct {
	commonDialect
}

func (p *mssql) Driver() string { return "mssql" }

func (p *mssql) BulidTopNSql(s *Scope, n int) (string, []interface{}) {
	wsql, args := s.buildWhere()
	sql := fmt.Sprintf("SELECT TOP %v * from %s %v", n, p.Quote(s.model.Name), wsql)
	return sql, args
}
