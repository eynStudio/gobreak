package orm

import (
	"fmt"
	"github.com/eynstudio/gobreak/db"
)

type mssql struct {
	commonDialect
}

func (p *mssql) Driver() string { return "mssql" }

func (p *mssql) BulidTopNSql(s *Scope, n int) (sa db.SqlArgs) {
	sa = s.buildWhere()
	sa.Sql = fmt.Sprintf("SELECT TOP %v * from %s %v", n, p.Quote(s.model.Name), sa.Sql)
	return
}
