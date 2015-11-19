package orm

import (
	"fmt"
)

type Dialect interface {
	Quote(key string) string
	Driver() string
}

func NewDialect(driver string) Dialect {
	var d Dialect
	switch driver {
	case "mysql":
		d = &mysql{}
	case "mssql":
		d = &mssql{}
	default:
		fmt.Printf("`%v` is not officially supported, running under compatibility mode.\n", driver)
		d = &commonDialect{}
	}
	return d
}

type commonDialect struct{}

func (commonDialect) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (p *commonDialect) Driver() string { return "common" }
