package dddd

import (
	"errors"
	. "github.com/eynstudio/gobreak/dddd/cmdbus"
	"reflect"
)

var (
	ErrAggregateAlreadySet = errors.New("aggregate is already set")
	ErrAggregateNotFound   = errors.New("no aggregate for command")
	cmdAggMap              = map[reflect.Type]reflect.Type{}
)

type aggCmdHandler struct {
}

func (p *aggCmdHandler) SetAgg(agg Agg) error {
	aggType := reflect.TypeOf(agg)
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

func (p *aggCmdHandler) Handle(cmd Cmd) error {

	return nil
}
