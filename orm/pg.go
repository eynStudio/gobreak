package orm

import (
	"encoding/json"
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
