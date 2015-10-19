package mgo

import (
	"errors"
	"time"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	. "github.com/eynstudio/gobreak/ddd"
	"github.com/eynstudio/gobreak/di"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrCouldNotDialDB         = errors.New("could not dial database")
	ErrNoDBSession            = errors.New("no database session")
	ErrCouldNotClearDB        = errors.New("could not clear database")
	ErrEventNotRegistered     = errors.New("event not registered")
	ErrModelNotSet            = errors.New("model not set")
	ErrCouldNotMarshalEvent   = errors.New("could not marshal event")
	ErrCouldNotUnmarshalEvent = errors.New("could not unmarshal event")
	ErrCouldNotLoadAggregate  = errors.New("could not load aggregate")
	ErrCouldNotSaveAggregate  = errors.New("could not save aggregate")
	ErrInvalidEvent           = errors.New("invalid event")
)

type MongoEventStore struct {
	eventBus     EventBus
	eventRepo    MgoRepo
	snapshotRepo MgoRepo
	factories    map[string]func() Event
}

func NewMongoEventStore(eventBus EventBus) (*MongoEventStore, error) {

	eventRepo := NewMgoRepo("DomainEvents", NewMongoAggregateRecord)
	di.Root.Apply(eventRepo)
	snapshotRepo := NewMgoRepo("DomainSnapshot", NewMongoAggregateRecord)
	di.Root.Apply(snapshotRepo)

	s := &MongoEventStore{
		eventBus:     eventBus,
		factories:    make(map[string]func() Event),
		eventRepo:    eventRepo,
		snapshotRepo: snapshotRepo,
	}

	return s, nil
}

type mongoAggregateRecord struct {
	Id      GUID                `bson:"_id"`
	Version int                 `Version`
	Events  []*mongoEventRecord `Events`
	// Type        string        `bson:"type"`
	// Snapshot    bson.Raw      `bson:"snapshot"`
}

func NewMongoAggregateRecord() T {
	return &mongoAggregateRecord{}
}

type mongoEventRecord struct {
	Type      string    `Type`
	Version   int       `Version`
	Timestamp time.Time `Timestamp`
	Event     Event     `bson:"-"`
	Data      bson.Raw  `Data`
}

func (s *MongoEventStore) Save(events []Event) error {
	if len(events) == 0 {
		return ErrNoEventsToAppend
	}

	sess := s.eventRepo.CopySession()
	defer sess.Close()
	c := s.eventRepo.C(sess)

	for _, event := range events {
		// Get an existing aggregate, if any.
		var existing []mongoAggregateRecord
		err := c.FindId(event.ID()).Select(bson.M{"Version": 1}).Limit(1).All(&existing)
		if err != nil || len(existing) > 1 {
			return ErrCouldNotLoadAggregate
		}

		var data []byte
		if data, err = bson.Marshal(event); err != nil {
			return ErrCouldNotMarshalEvent
		}

		r := &mongoEventRecord{
			Type:      event.EventType(),
			Version:   1,
			Timestamp: time.Now(),
			Data:      bson.Raw{3, data},
		}

		if len(existing) == 0 {
			aggregate := mongoAggregateRecord{
				Id:      event.ID(),
				Version: 1,
				Events:  []*mongoEventRecord{r},
			}

			if err := c.Insert(aggregate); err != nil {
				return ErrCouldNotSaveAggregate
			}
		} else {
			r.Version = existing[0].Version + 1

			// Increment aggregate version on insert of new event record, and
			// only insert if version of aggregate is matching (ie not changed
			// since the query above).
			err = c.Update(
				bson.M{
					"_id":     event.ID(),
					"Version": existing[0].Version,
				},
				bson.M{
					"$push": bson.M{"Events": r},
					"$inc":  bson.M{"Version": 1},
				},
			)
			if err != nil {
				return ErrCouldNotSaveAggregate
			}
		}

		if s.eventBus != nil {
			s.eventBus.PublishEvent(event)
		}
	}

	return nil
}

func (s *MongoEventStore) SaveSnapshot(agg Aggregate) error {
	s.snapshotRepo.Save(agg.ID(),agg.GetSnapshot())
	return nil
}

func (s *MongoEventStore) LoadSnapshot(agg Aggregate) error{
	sess :=s.snapshotRepo.CopySession()
	defer sess.Close()
	s.snapshotRepo.C(sess).FindId(agg.ID()).One(agg.GetSnapshot())
	return nil
}
func (s *MongoEventStore) Load(id GUID) ([]Event, error) {
	sess := s.eventRepo.CopySession()
	defer sess.Close()
	c := s.eventRepo.C(sess)

	var aggregate mongoAggregateRecord
	err := c.FindId(id).One(&aggregate)
	if err != nil {
		return nil, ErrNoEventsFound
	}

	events := make([]Event, len(aggregate.Events))
	for i, record := range aggregate.Events {
		f, ok := s.factories[record.Type]
		if !ok {
			return nil, ErrEventNotRegistered
		}

		event := f()
		if err := record.Data.Unmarshal(event); err != nil {
			return nil, ErrCouldNotUnmarshalEvent
		}
		if events[i], ok = event.(Event); !ok {
			return nil, ErrInvalidEvent
		}

		record.Event = events[i]
		record.Data = bson.Raw{}
	}

	return events, nil
}

// RegisterEventType registers an event factory for a event type. The factory is
// used to create concrete event types when loading from the database.
//
// An example would be:
//     eventStore.RegisterEventType(&MyEvent{}, func() Event { return &MyEvent{} })
func (s *MongoEventStore) RegisterEventType(event Event, factory func() Event) error {
	if _, ok := s.factories[event.EventType()]; ok {
		return ErrHandlerExist
	}

	s.factories[event.EventType()] = factory
	return nil
}

func (s *MongoEventStore) Clear() error {
	sess := s.eventRepo.CopySession()
	defer sess.Close()
	c := s.eventRepo.C(sess)
	if err := c.DropCollection(); err != nil {
		return ErrCouldNotClearDB
	}
	return nil
}
