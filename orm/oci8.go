package orm

type oci8 struct {
	commonDialect
}

func (p *oci8) Driver() string { return "oci8" }
