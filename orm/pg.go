package orm

type pg struct {
	commonDialect
}

func (p *pg) Driver() string {
	return "postgres"
}
