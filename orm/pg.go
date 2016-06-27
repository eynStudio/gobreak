package orm

import (
	"encoding/json"
	"fmt"

	"github.com/eynstudio/gobreak/db"
)

type pg struct {
	commonDialect
}

func (p *pg) Driver() string {
	return "postgres"
}

type Jsonb struct {
	Id   string          `Id`
	Json json.RawMessage `Json`
}

func (p *pg) BulidTopNSql(s *Scope, n int) (sa db.SqlArgs) {
	sa = s.buildWhere()
	sa.Sql = fmt.Sprintf("SELECT * from %s %v limit %d", p.Quote(s.model.Name), sa.Sql, n)
	return
}
