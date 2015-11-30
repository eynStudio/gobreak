package orm

import (
	"fmt"
)

type mysql struct {
	commonDialect
}

func (p *mysql) Driver() string { return "mysql" }

func (mysql) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}