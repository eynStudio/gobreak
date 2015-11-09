package orm

import (
	"fmt"
)

import (
	. "github.com/eynstudio/gobreak"
)

type Where struct {
	model *model
	Dialect
}

func NewWhere(model *model, dialect Dialect) *Where {
	w := &Where{model, dialect}

	return w
}

func (p *Where) where(data T, params ...interface{}) (string, []interface{}) {
	var args []interface{}
	sql := ""
	if len(params) == 0 {
		return sql,args
	}
	if len(params) > 0 {
		sql = fmt.Sprintf("WHERE %v", params[0])
	}
	if len(params) > 1 {
		args = params[1:]
	}

	return sql, args
}

//func (p *Where) whereId(data T, idval interface{}) (string, []interface{}) {
//	id := p.model.Id()
//	var idval interface{}
//	if len(idVals) == 0 {
//		idval = p.model.IdVal(data)
//	} else {
//		idval = idVals[0]
//	}
//	sql := fmt.Sprintf("WHERE %v=?", id)
//	return sql, []interface{}{idval}
//}
