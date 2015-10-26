package ddd

import (
	"errors"
	. "github.com/eynstudio/gobreak"
	"reflect"
)

var (
	ErrAggregateAlreadyRegistered = errors.New("aggregate is already registered")
	ErrAggregateNotRegistered     = errors.New("aggregate is not registered")
)

type DomainRepo struct {
	EventStore EventStore `di`
	callbacks  map[reflect.Type]func(GUID) Aggregate
}

func NewDomainRepo() *DomainRepo {
	return &DomainRepo{
		callbacks: make(map[reflect.Type]func(GUID) Aggregate),
	}
}

func (p *DomainRepo) RegisterAggregate(aggregate Aggregate, callback func(GUID) Aggregate) error {
	aggType := reflect.TypeOf(aggregate)
	if _, ok := p.callbacks[aggType]; ok {
		return ErrAggregateAlreadyRegistered
	}

	p.callbacks[aggType] = callback
	return nil
}

func (p *DomainRepo) Load(aggregateType reflect.Type, id GUID) (Aggregate, error) {
	if f, ok := p.callbacks[aggregateType]; ok {
		return p.EventStore.Load(f(id))
	} else {
		return nil, ErrAggregateNotRegistered
	}
}

func (p *DomainRepo) Save(aggregate Aggregate) error {
	if aggregate.HasUncommittedEvents() {
		return p.EventStore.Save(aggregate)
	}
	return nil
}
