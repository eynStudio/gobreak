package orm

import (
	"strings"
	. "github.com/eynstudio/gobreak"
	"fmt"
)

type sqlBuilder struct{
	model *model
	dialect Dialect
}

func newSqlBuilder(m *model,dialect Dialect) *sqlBuilder{
	return &sqlBuilder{
		model:m,
		dialect:dialect,
	}
}
func (p *sqlBuilder) buildSelect() string{
	sql:=fmt.Sprintf("select * from %s",p.dialect.Quote(p.model.Name))	
	return sql
}

func (p *sqlBuilder) buildSelect1() string{
	sql:=fmt.Sprintf("select top 1 * from %s",p.dialect.Quote(p.model.Name))	
	return sql
}

func (p *sqlBuilder) buildInsert(obj T) (string,[]interface{}){
	var cols []string
	var params []string
	var args []interface{}
	m:=p.model.Obj2Map(obj)
	for k,v:=range m{
		cols=append(cols,p.dialect.Quote(k))
		params=append(params,"?")
		args=append(args,v)
	}
	sql:=fmt.Sprintf("insert into %s (%v) values (%v)",
		p.dialect.Quote(p.model.Name),
		strings.Join(cols,","),
		strings.Join( params,","),
	)	
	return sql,args
}

func (p *sqlBuilder) buildUpdate(obj T) (string,[]interface{}){
	id:=p.model.Id()
	idval:=p.model.IdVal(obj)	
	
	var cols []string
	var args []interface{}
	m:=p.model.Obj2Map(obj)
	for k,v:=range m{
		cols=append(cols,p.dialect.Quote(k)+"=?")
		args=append(args,v)
	}
	args=append(args,idval)
	
	sql:=fmt.Sprintf("UPDATE %s SET %v where %v=?",
		p.dialect.Quote(p.model.Name),
		strings.Join(cols,","),
		p.dialect.Quote(id),
	)	
	return sql,args
}