package ddd

import (
	"errors"
	. "github.com/eynstudio/gobreak"
)

var (
	ErrAggregateAlreadyRegistered = errors.New("aggregate is already registered")
	ErrAggregateNotRegistered     = errors.New("aggregate is not registered")
)

type DomainRepo struct {
	EventStore EventStore `di`
	callbacks  map[string]func(GUID) Aggregate
}

func NewDomainRepo() *DomainRepo {
	return &DomainRepo{
		callbacks: make(map[string]func(GUID) Aggregate),
	}
}

func (p *DomainRepo) RegisterAggregate(aggregate Aggregate, callback func(GUID) Aggregate) error {
	if _, ok := p.callbacks[aggregate.AggType()]; ok {
		return ErrAggregateAlreadyRegistered
	}

	p.callbacks[aggregate.AggType()] = callback

	return nil
}

func (p *DomainRepo) Load(aggregateType string, id GUID) (Aggregate, error) {
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
