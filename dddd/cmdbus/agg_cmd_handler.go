package cmdbus

import (
	"errors"
	"reflect"

	. "github.com/eynstudio/gobreak/dddd/ddd"
	"github.com/eynstudio/gobreak/dddd/store"
	"github.com/eynstudio/gobreak/di"
)

var (
	ErrAggregateAlreadySet = errors.New("aggregate is already set")
	ErrAggregateNotFound   = errors.New("no aggregate for command")
	cmdAggMap              = map[reflect.Type]reflect.Type{}
)

type aggCmdHandler struct {
}

func SetAgg(agg Agg) error {
	aggType := reflect.TypeOf(agg)
	aggType = di.Ptr2Elem(aggType)
	for _, c := range agg.RegistedCmds() {
		cmdType := reflect.TypeOf(c)
		if _, ok := cmdAggMap[cmdType]; ok {
			return ErrAggregateAlreadySet
		}
		cmdAggMap[cmdType] = aggType
	}
	return nil
}

func (p *aggCmdHandler) CanHandle(cmd Cmd) bool {
	_, ok := cmdAggMap[reflect.TypeOf(cmd)]
	return ok
}

func (p *aggCmdHandler) Handle(cmd Cmd) (err error) {
	var agg Agg
	if aggType, ok := cmdAggMap[reflect.TypeOf(cmd)]; !ok {
		return ErrAggregateNotFound
	} else if agg, err = store.Load(aggType, cmd.ID()); err == nil {
		if err = agg.HandleCmd(cmd); err == nil {
			err = store.Save(agg)
		}
	}
	return
}
