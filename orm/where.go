package orm

import (
	"fmt"
)

import(
	. "github.com/eynstudio/gobreak"
)
type Where struct{
		model *model
	Dialect	
}

func NewWhere(model *model,dialect Dialect) *Where{
	w:=&Where{model, dialect}
	
	return w
}

func (p*Where) w(data T,params ...interface{}) (string,[]interface{}){
	if len(params)==0{
		return p.whereId(data)
	}

	var args []interface{}
	sql:=""
	
	if len(params)>0 {
		sql=fmt.Sprintf("WHERE %v",params[0])
	}
	if len(params)>1{
		args=params[1:]
	}
	
	return sql,args
}

func (p*Where) whereId(data T) (string,[]interface{}){
	id:=p.model.Id()
	idval:=p.model.IdVal(data)
	sql:=fmt.Sprintf("WHERE %v=?",id)
	return sql,[]interface{}{idval}
}