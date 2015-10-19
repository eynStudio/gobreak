package ddd

import (
	"errors"
	. "github.com/eynstudio/gobreak"
	"time"
)

var (
	ErrNoEventsToAppend    = errors.New("no events to append")
	ErrNoEventsFound       = errors.New("could not find events")
	ErrNoEventStoreDefined = errors.New("no event store defined")
)

type EventStore interface {
	Save([]Event) error
	SaveSnapshot(agg Aggregate) error
	Load(id GUID) ([]Event, error)
		LoadSnapshot(agg Aggregate) error
}

type AggregateRecord interface {
	ID() GUID
	Version() int
	EventRecords() []EventRecord
}

type EventRecord interface {
	Type() string
	Version() int
	Events() []Event
}

type MemoryEventStore struct {
	eventBus         EventBus
	aggregateRecords map[GUID]*memoryAggregateRecord
}

func NewMemoryEventStore(eventBus EventBus) *MemoryEventStore {
	s := &MemoryEventStore{
		eventBus:         eventBus,
		aggregateRecords: make(map[GUID]*memoryAggregateRecord),
	}
	return s
}

func (s *MemoryEventStore) Save(events []Event) error {
	if len(events) == 0 {
		return ErrNoEventsToAppend
	}

	for _, event := range events {
		r := &memoryEventRecord{
			eventType: event.EventType(),
			timestamp: time.Now(),
			event:     event,
		}

		if a, ok := s.aggregateRecords[event.ID()]; ok {
			a.version++
			r.version = a.version
			a.events = append(a.events, r)
		} else {
			s.aggregateRecords[event.ID()] = &memoryAggregateRecord{
				aggregateID: event.ID(),
				version:     0,
				events:      []*memoryEventRecord{r},
			}
		}

		if s.eventBus != nil {
			s.eventBus.PublishEvent(event)
		}
	}

	return nil
}

func (s *MemoryEventStore) Load(id GUID) ([]Event, error) {
	if a, ok := s.aggregateRecords[id]; ok {
		events := make([]Event, len(a.events))
		for i, r := range a.events {
			events[i] = r.event
		}
		return events, nil
	}

	return nil, ErrNoEventsFound
}

type memoryAggregateRecord struct {
	aggregateID GUID
	version     int
	events      []*memoryEventRecord
}

type memoryEventRecord struct {
	eventType string
	version   int
	timestamp time.Time
	event     Event
}
