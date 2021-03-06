package orm

import (
	"database/sql"
	"fmt"

	"github.com/eynstudio/gobreak/db"
	"github.com/eynstudio/gobreak/db/meta"
)

type Dialect interface {
	Quote(key string) string
	Driver() string
	BulidTopNSql(s *Scope, n int) db.SqlArgs
	LoadMeta(*sql.DB) *meta.MetaDb
}

func NewDialect(driver string) Dialect {
	var d Dialect
	switch driver {
	case "mysql":
		d = &mysql{}
	case "postgres":
		d = &pg{}
	case "mssql":
		d = &mssql{}
	case "oci8":
		d = &oci8{}
	default:
		fmt.Printf("`%v` is not officially supported, running under compatibility mode.\n", driver)
		d = &commonDialect{}
	}
	return d
}

type commonDialect struct{}

func (commonDialect) Quote(key string) string                         { return fmt.Sprintf(`"%s"`, key) }
func (p *commonDialect) Driver() string                               { return "common" }
func (p *commonDialect) BulidTopNSql(s *Scope, n int) (sa db.SqlArgs) { return }
func (p *commonDialect) LoadMeta(*sql.DB) *meta.MetaDb                { return nil }
