package ddd

import (
	"errors"
	"reflect"
)

var (
	ErrAggregateAlreadySet = errors.New("aggregate is already set")
	ErrAggregateNotFound   = errors.New("no aggregate for command")
)

type CmdFieldError struct {
	Field string
}

func (p CmdFieldError) Error() string {
	return "missing field: " + p.Field
}

type AggCmdHandler interface {
	CmdHandler
	SetAggregate(agg Aggregate) error
}

type aggCmdHandler struct {
	DomainRepo `di`
	cmdAggMap  map[reflect.Type]reflect.Type
}

func NewAggCmdHandler() AggCmdHandler {
	return &aggCmdHandler{
		cmdAggMap: make(map[reflect.Type]reflect.Type),
	}
}

func (p *aggCmdHandler) SetAggregate(agg Aggregate) error {
	aggType := reflect.TypeOf(agg)
	for _, c := range agg.RegistedCmds() {
		cmdType := reflect.TypeOf(c)
		if _, ok := p.cmdAggMap[cmdType]; ok {
			return ErrAggregateAlreadySet
		}
		p.cmdAggMap[cmdType] = aggType
	}
	return nil
}

func (p *aggCmdHandler) CanHandleCmd(cmd Cmd) bool {
	_, ok := p.cmdAggMap[reflect.TypeOf(cmd)]
	return ok
}

func (p *aggCmdHandler) HandleCmd(cmd Cmd) error {
	var err error
	var aggregateType reflect.Type
	var ok bool
	if aggregateType, ok = p.cmdAggMap[reflect.TypeOf(cmd)]; !ok {
		return ErrAggregateNotFound
	}

	var aggregate Aggregate
	if aggregate, err = p.DomainRepo.Load(aggregateType, cmd.ID()); err != nil {
		return err
	}

	if err = aggregate.HandleCmd(cmd); err != nil {
		return err
	}

	if err = p.DomainRepo.Save(aggregate); err != nil {
		return err
	}

	return nil
}
