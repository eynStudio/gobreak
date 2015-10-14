package mgo

import (
	"errors"
	. "github.com/eynstudio/gobreak/ddd"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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
	eventBus  EventBus
	session   *mgo.Session
	db        string
	factories map[string]func() Event
}

func NewMongoEventStore(eventBus EventBus, url, database string) (*MongoEventStore, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, ErrCouldNotDialDB
	}

	session.SetMode(mgo.Strong, true)
	session.SetSafe(&mgo.Safe{W: 1})

	return NewMongoEventStoreWithSession(eventBus, session, database)
}

func NewMongoEventStoreWithSession(eventBus EventBus, session *mgo.Session, database string) (*MongoEventStore, error) {
	if session == nil {
		return nil, ErrNoDBSession
	}

	s := &MongoEventStore{
		eventBus:  eventBus,
		factories: make(map[string]func() Event),
		session:   session,
		db:        database,
	}

	return s, nil
}

type mongoAggregateRecord struct {
	AggregateID GUID                `bson:"_id"`
	Version     int                 `Version`
	Events      []*mongoEventRecord `Events`
	// Type        string        `bson:"type"`
	// Snapshot    bson.Raw      `bson:"snapshot"`
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

	sess := s.session.Copy()
	defer sess.Close()

	for _, event := range events {
		// Get an existing aggregate, if any.
		var existing []mongoAggregateRecord
		err := sess.DB(s.db).C("events").FindId(event.ID()).
			Select(bson.M{"version": 1}).Limit(1).All(&existing)
		if err != nil || len(existing) > 1 {
			return ErrCouldNotLoadAggregate
		}

		// Marshal event data.
		var data []byte
		if data, err = bson.Marshal(event); err != nil {
			return ErrCouldNotMarshalEvent
		}

		// Create the event record with timestamp.
		r := &mongoEventRecord{
			Type:      event.EventType(),
			Version:   1,
			Timestamp: time.Now(),
			Data:      bson.Raw{3, data},
		}

		// Either insert a new aggregate or append to an existing.
		if len(existing) == 0 {
			aggregate := mongoAggregateRecord{
				AggregateID: event.ID(),
				Version:     1,
				Events:      []*mongoEventRecord{r},
			}

			if err := sess.DB(s.db).C("events").Insert(aggregate); err != nil {
				return ErrCouldNotSaveAggregate
			}
		} else {
			// Increment record version before inserting.
			r.Version = existing[0].Version + 1

			// Increment aggregate version on insert of new event record, and
			// only insert if version of aggregate is matching (ie not changed
			// since the query above).
			err = sess.DB(s.db).C("events").Update(
				bson.M{
					"_id":     event.ID(),
					"version": existing[0].Version,
				},
				bson.M{
					"$push": bson.M{"events": r},
					"$inc":  bson.M{"version": 1},
				},
			)
			if err != nil {
				return ErrCouldNotSaveAggregate
			}
		}

		// Publish event on the bus.
		if s.eventBus != nil {
			s.eventBus.PublishEvent(event)
		}
	}

	return nil
}

func (s *MongoEventStore) Load(id GUID) ([]Event, error) {
	sess := s.session.Copy()
	defer sess.Close()

	var aggregate mongoAggregateRecord
	err := sess.DB(s.db).C("events").FindId(id).One(&aggregate)
	if err != nil {
		return nil, ErrNoEventsFound
	}

	events := make([]Event, len(aggregate.Events))
	for i, record := range aggregate.Events {
		// Get the registered factory function for creating events.
		f, ok := s.factories[record.Type]
		if !ok {
			return nil, ErrEventNotRegistered
		}

		// Manually decode the raw BSON event.
		event := f()
		if err := record.Data.Unmarshal(event); err != nil {
			return nil, ErrCouldNotUnmarshalEvent
		}
		if events[i], ok = event.(Event); !ok {
			return nil, ErrInvalidEvent
		}

		// Set conrcete event and zero out the decoded event.
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

// SetDB sets the database session.
func (s *MongoEventStore) SetDB(db string) {
	s.db = db
}

// Clear clears the event storge.
func (s *MongoEventStore) Clear() error {
	if err := s.session.DB(s.db).C("events").DropCollection(); err != nil {
		return ErrCouldNotClearDB
	}
	return nil
}

// Close closes the database session.
func (s *MongoEventStore) Close() {
	s.session.Close()
}
