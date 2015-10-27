package ddd

import (
	"errors"
	"reflect"
	"time"
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
	err := p.checkCmd(cmd)
	if err != nil {
		return err
	}

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

func (p *aggCmdHandler) checkCmd(cmd Cmd) error {
	rv := reflect.Indirect(reflect.ValueOf(cmd))
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue // Skip private field.
		}

		tag := field.Tag.Get("eh")
		if tag == "optional" {
			continue // Optional field.
		}

		if isZero(rv.Field(i)) {
			return CmdFieldError{field.Name}
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Struct:
		// Special case to get zero values by method.
		switch obj := v.Interface().(type) {
		case time.Time:
			return obj.IsZero()
		}

		// Check public fields for zero values.
		z := true
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).PkgPath != "" {
				continue // Skip private fields.
			}
			z = z && isZero(v.Field(i))
		}
		return z
	}

	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
