package db

type M map[string]interface{}

type SqlArgs struct {
	Sql  string
	Args []interface{}
}

func (p *SqlArgs) AddArgs(a ...interface{}) { p.Args = append(p.Args, a...) }
