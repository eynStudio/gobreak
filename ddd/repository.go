package ddd

import (
	"errors"
)

var (
	ErrNilEventStore              = errors.New("event store is nil")
	ErrAggregateAlreadyRegistered = errors.New("aggregate is already registered")
	ErrAggregateNotRegistered     = errors.New("aggregate is not registered")
)

type Repository interface {
	Load(string, GUID) (Aggregate, error)
	Save(Aggregate) error
}

type CallbackRepository struct {
	eventStore EventStore
	callbacks  map[string]func(GUID) Aggregate
}

func NewCallbackRepository(eventStore EventStore) (*CallbackRepository, error) {
	if eventStore == nil {
		return nil, ErrNilEventStore
	}

	d := &CallbackRepository{
		eventStore: eventStore,
		callbacks:  make(map[string]func(GUID) Aggregate),
	}
	return d, nil
}

func (r *CallbackRepository) RegisterAggregate(aggregate Aggregate, callback func(GUID) Aggregate) error {
	if _, ok := r.callbacks[aggregate.AggType()]; ok {
		return ErrAggregateAlreadyRegistered
	}

	r.callbacks[aggregate.AggType()] = callback

	return nil
}

func (r *CallbackRepository) Load(aggregateType string, id GUID) (Aggregate, error) {
	f, ok := r.callbacks[aggregateType]
	if !ok {
		return nil, ErrAggregateNotRegistered
	}

	aggregate := f(id)
	events, _ := r.eventStore.Load(aggregate.ID())
	for _, event := range events {
		aggregate.ApplyEvent(event)
		aggregate.IncrementVersion()
	}

	return aggregate, nil
}

func (r *CallbackRepository) Save(aggregate Aggregate) error {
	resultEvents := aggregate.GetUncommittedEvents()

	if len(resultEvents) > 0 {
		err := r.eventStore.Save(resultEvents)
		if err != nil {
			return err
		}
	}

	aggregate.ClearUncommittedEvents()
	return nil
}
