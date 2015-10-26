package ddd

import (
	"errors"
	"reflect"
	"time"
)

var (
	ErrNilRepository       = errors.New("repository is nil")
	ErrAggregateAlreadySet = errors.New("aggregate is already set")
	ErrAggregateNotFound   = errors.New("no aggregate for command")
)

type CommandFieldError struct {
	Field string
}

func (p CommandFieldError) Error() string {
	return "missing field: " + p.Field
}

type AggregateCommandHandler struct {
	repository Repository
	cmdAggMap  map[reflect.Type]reflect.Type
}

func NewAggregateCommandHandler(repository Repository) (*AggregateCommandHandler, error) {
	if repository == nil {
		return nil, ErrNilRepository
	}

	h := &AggregateCommandHandler{
		repository: repository,
		cmdAggMap:  make(map[reflect.Type]reflect.Type),
	}
	return h, nil
}

func (p *AggregateCommandHandler) SetAggregateCmd(aggregate Aggregate, cmd Cmd) error {
	if _, ok := p.cmdAggMap[reflect.TypeOf(cmd)]; ok {
		return ErrAggregateAlreadySet
	}

	p.cmdAggMap[reflect.TypeOf(cmd)] = reflect.TypeOf(aggregate)
	return nil
}

func (p *AggregateCommandHandler) SetAggregateCmds(aggregate Aggregate, cmds ...Cmd) error {
	var err error
	for _, c := range cmds {
		err = p.SetAggregateCmd(aggregate, c)
	}
	return err
}

func (p *AggregateCommandHandler) CanHandleCmd(cmd Cmd) bool {
	_, ok := p.cmdAggMap[reflect.TypeOf(cmd)]
	return ok
}

func (p *AggregateCommandHandler) HandleCmd(cmd Cmd) error {
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
	if aggregate, err = p.repository.Load(aggregateType, cmd.ID()); err != nil {
		return err
	}

	if err = aggregate.HandleCmd(cmd); err != nil {
		return err
	}

	if err = p.repository.Save(aggregate); err != nil {
		return err
	}

	return nil
}

func (p *AggregateCommandHandler) checkCmd(cmd Cmd) error {
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
			return CommandFieldError{field.Name}
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
