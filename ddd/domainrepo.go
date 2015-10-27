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

type DomainRepo interface {
	Load(reflect.Type, GUID) (Aggregate, error)
	Save(Aggregate) error
	RegisterAggregate(aggregate Aggregate, callback func(GUID) Aggregate) error
}

type domainRepo struct {
	EventStore EventStore `di`
	callbacks  map[reflect.Type]func(GUID) Aggregate
}

func NewDomainRepo() DomainRepo {
	return &domainRepo{
		callbacks: make(map[reflect.Type]func(GUID) Aggregate),
	}
}

func (p *domainRepo) RegisterAggregate(aggregate Aggregate, callback func(GUID) Aggregate) error {
	aggType := reflect.TypeOf(aggregate)
	if _, ok := p.callbacks[aggType]; ok {
		return ErrAggregateAlreadyRegistered
	}

	p.callbacks[aggType] = callback
	return nil
}

func (p *domainRepo) Load(aggregateType reflect.Type, id GUID) (Aggregate, error) {
	if f, ok := p.callbacks[aggregateType]; ok {
		return p.EventStore.Load(f(id))
	} else {
		return nil, ErrAggregateNotRegistered
	}
}

func (p *domainRepo) Save(aggregate Aggregate) error {
	if aggregate.HasUncommittedEvents() {
		return p.EventStore.Save(aggregate)
	}
	return nil
}
