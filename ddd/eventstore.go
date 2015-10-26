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
	Save(agg Aggregate) error
	Load(agg Aggregate) (Aggregate, error)
}

type EventRecord interface {
	Type() string
	Version() int
	Events() []Event
}

type MemoryEventStore struct {
	eventBus         EventBus
	aggregateRecords map[GUID]Aggregate
	eventRecords     []memoryEventRecord
}

func NewMemoryEventStore(eventBus EventBus) *MemoryEventStore {
	s := &MemoryEventStore{
		eventBus:         eventBus,
		aggregateRecords: make(map[GUID]Aggregate),
	}
	return s
}

func (s *MemoryEventStore) Save(agg Aggregate) error {
	s.aggregateRecords[agg.ID()] = agg

	if !agg.HasUncommittedEvents() {
		return nil
	}

	events := agg.GetUncommittedEvents()
	for _, event := range events {
		r := &memoryEventRecord{
			eventType: event.EventType(),
			timestamp: time.Now(),
			event:     event,
		}
		s.eventRecords = append(s.eventRecords, *r)
		if s.eventBus != nil {
			s.eventBus.PublishEvent(event)
		}
	}
	agg.ClearUncommittedEvents()
	return nil
}

func (s *MemoryEventStore) Load(agg Aggregate) (Aggregate, error) {
	if a, ok := s.aggregateRecords[agg.ID()]; ok {
		return a, nil
	}
	return agg, nil
}

type memoryEventRecord struct {
	eventType string
	timestamp time.Time
	event     Event
}
