package orm

import (
	"fmt"
	"github.com/eynstudio/gobreak/db"
)

type oci8 struct {
	commonDialect
}

func (p *oci8) Driver() string { return "oci8" }

func (p *oci8) BulidTopNSql(s *Scope, n int) (sa db.SqlArgs) {
	sa = s.buildWhere()
	if sa.Sql == "" {
		sa.Sql = fmt.Sprintf("Where ROWNUM <= %d", n)
	} else {
		sa.Sql = sa.Sql + fmt.Sprintf(" and (ROWNUM <= %d)", n)
	}
	sa.Sql = fmt.Sprintf("SELECT * from %s %v", p.Quote(s.model.Name), sa.Sql)
	return
}
