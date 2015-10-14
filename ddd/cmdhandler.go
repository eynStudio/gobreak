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

func (c CommandFieldError) Error() string {
	return "missing field: " + c.Field
}

// AggregateCommandHandler dispatches commands to registered aggregates.
//
// The dispatch process is as follows:
// 1. The handler receives a command
// 2. An aggregate is created or rebuilt from previous events by the repository
// 3. The aggregate's command handler is called
// 4. The aggregate stores events in response to the command
// 5. The new events are stored in the event store by the repository
// 6. The events are published to the event bus when stored by the event store
type AggregateCommandHandler struct {
	repository Repository
	aggregates map[string]string
}

func NewAggregateCommandHandler(repository Repository) (*AggregateCommandHandler, error) {
	if repository == nil {
		return nil, ErrNilRepository
	}

	h := &AggregateCommandHandler{
		repository: repository,
		aggregates: make(map[string]string),
	}
	return h, nil
}

func (h *AggregateCommandHandler) SetAggregate(aggregate Aggregate, cmd Cmd) error {
	if _, ok := h.aggregates[cmd.CmdType()]; ok {
		return ErrAggregateAlreadySet
	}

	h.aggregates[cmd.CmdType()] = aggregate.AggType()
	return nil
}

func (h *AggregateCommandHandler) HandleCmd(cmd Cmd) error {
	err := h.checkCmd(cmd)
	if err != nil {
		return err
	}

	var aggregateType string
	var ok bool
	if aggregateType, ok = h.aggregates[cmd.CmdType()]; !ok {
		return ErrAggregateNotFound
	}

	var aggregate Aggregate
	if aggregate, err = h.repository.Load(aggregateType, cmd.ID()); err != nil {
		return err
	}

	if err = aggregate.HandleCmd(cmd); err != nil {
		return err
	}

	if err = h.repository.Save(aggregate); err != nil {
		return err
	}

	return nil
}

func (h *AggregateCommandHandler) checkCmd(cmd Cmd) error {
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
