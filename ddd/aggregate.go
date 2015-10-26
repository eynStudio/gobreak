package ddd

import (
	. "github.com/eynstudio/gobreak"
)

type Aggregate interface {
	ID() GUID
	Version() int
	IncrementVersion()
	HandleCmd(Cmd) error
	ApplyEvent(events Event)
	GetSnapshot() T
	StoreEvent(Event)
	GetUncommittedEvents() []Event
	ClearUncommittedEvents()
	HasUncommittedEvents() bool
}

type AggregateBase struct {
	id                GUID
	version           int
	uncommittedEvents []Event
}

func NewAggregateBase(id GUID) *AggregateBase {
	return &AggregateBase{
		id:                id,
		uncommittedEvents: []Event{},
	}
}

func (a *AggregateBase) ID() GUID {
	return a.id
}

func (a *AggregateBase) Version() int {
	return a.version
}

func (a *AggregateBase) IncrementVersion() {
	a.version++
}

func (a *AggregateBase) StoreEvent(event Event) {
	a.uncommittedEvents = append(a.uncommittedEvents, event)
}

func (a *AggregateBase) GetUncommittedEvents() []Event {
	return a.uncommittedEvents
}

func (a *AggregateBase) HasUncommittedEvents() bool {
	return len(a.uncommittedEvents) > 0
}

func (a *AggregateBase) ClearUncommittedEvents() {
	a.uncommittedEvents = []Event{}
}
