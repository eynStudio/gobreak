package orm

import (
	"fmt"
)

type oci8 struct {
	commonDialect
}

func (p *oci8) Driver() string { return "oci8" }

func (p *oci8) BulidTopNSql(s *Scope, n int) (string, []interface{}) {
	wsql, args := s.buildWhere()
	if wsql == "" {
		wsql = fmt.Sprintf("Where ROWNUM <= %d", n)
	} else {
		wsql = wsql + fmt.Sprintf(" and (ROWNUM <= %d)", n)
	}
	sql := fmt.Sprintf("SELECT * from %s %v", p.Quote(s.model.Name), wsql)
	return sql, args
}
