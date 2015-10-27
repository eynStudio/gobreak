package mgo

import (
	"errors"
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	. "github.com/eynstudio/gobreak/ddd"
	"github.com/eynstudio/gobreak/di"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"time"
)

var (
	ErrCouldNotClearDB      = errors.New("could not clear database")
	ErrCouldNotMarshalEvent = errors.New("could not marshal event")
)

type MongoEventStore struct {
	EventBus     `di`
	eventRepo    MgoRepo
	snapshotRepo MgoRepo
	factories    map[string]func() Event
}

func NewMongoEventStore() (*MongoEventStore, error) {

	eventRepo := NewMgoRepo("SysAggEvent", nil)
	di.Root.Apply(eventRepo)
	snapshotRepo := NewMgoRepo("SysAggSnapshot", nil)
	di.Root.Apply(snapshotRepo)

	s := &MongoEventStore{
		factories:    make(map[string]func() Event),
		eventRepo:    eventRepo,
		snapshotRepo: snapshotRepo,
	}

	return s, nil
}

type sysDomainEvent struct {
	AggId GUID      `Id`
	Time  time.Time `Time`
	Type  string    `Type`
	Data  bson.Raw  `Data`
}

func (s *MongoEventStore) Save(agg Aggregate) error {
	events := agg.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	sess := s.eventRepo.CopySession()
	defer sess.Close()
	c := s.eventRepo.C(sess)

	for _, event := range events {

		var data []byte
		var err error
		if data, err = bson.Marshal(event); err != nil {
			return ErrCouldNotMarshalEvent
		}

		r := &sysDomainEvent{
			AggId: event.ID(),
			Type:  reflect.TypeOf(event).String(),
			Time:  time.Now(),
			Data:  bson.Raw{3, data},
		}

		c.Insert(r)

		if s.EventBus != nil {
			s.EventBus.PublishEvent(event)
		}
	}

	s.snapshotRepo.Save(agg.ID(), agg.GetSnapshot())
	agg.ClearUncommittedEvents()
	return nil
}

func (s *MongoEventStore) Load(agg Aggregate) (Aggregate, error) {
	sess := s.snapshotRepo.CopySession()
	defer sess.Close()
	s.snapshotRepo.C(sess).FindId(agg.ID()).One(agg.GetSnapshot())
	return agg, nil
}

func (s *MongoEventStore) Clear() error {
	sess := s.eventRepo.CopySession()
	defer sess.Close()
	c := s.eventRepo.C(sess)
	if err := c.DropCollection(); err != nil {
		return ErrCouldNotClearDB
	}
	c = s.snapshotRepo.C(sess)
	if err := c.DropCollection(); err != nil {
		return ErrCouldNotClearDB
	}
	return nil
}
