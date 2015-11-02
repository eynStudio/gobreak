package orm

import(
	"fmt"
)
type Dialect interface {

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