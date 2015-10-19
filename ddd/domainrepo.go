package ddd

import (
	"errors"
	. "github.com/eynstudio/gobreak"
)

var (
	ErrNilEventStore              = errors.New("event store is nil")
	ErrAggregateAlreadyRegistered = errors.New("aggregate is already registered")
	ErrAggregateNotRegistered     = errors.New("aggregate is not registered")
)

type DomainRepo struct {
	EventStore EventStore 	`di`
	callbacks  map[string]func(GUID) Aggregate
}

func NewDomainRepo() *DomainRepo{
	return &DomainRepo{
		callbacks:  make(map[string]func(GUID) Aggregate),
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
	f, ok := p.callbacks[aggregateType]
	if !ok {
		return nil, ErrAggregateNotRegistered
	}

	aggregate := f(id)
	p.EventStore.LoadSnapshot(aggregate)
//	events, _ := p.EventStore.Load(aggregate.ID())
//	for _, event := range events {
//		aggregate.ApplyEvent(event)
//		aggregate.IncrementVersion()
//	}

	return aggregate, nil
}

func (p *DomainRepo) Save(aggregate Aggregate) error {
	resultEvents := aggregate.GetUncommittedEvents()

	if len(resultEvents) > 0 {
		
		err := p.EventStore.Save(resultEvents)
		if err != nil {
			return err
		}
		p.EventStore.SaveSnapshot(aggregate)
	}

	aggregate.ClearUncommittedEvents()
	return nil
}
