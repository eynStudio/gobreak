package orm

type mysql struct {
	commonDialect
}

func (p *mysql) Driver() string { return "mysql" }
