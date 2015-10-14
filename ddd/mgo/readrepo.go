package mgo

import (
	. "github.com/eynstudio/gobreak/ddd"
	"gopkg.in/mgo.v2"
)

type MongoReadRepository struct {
	session    *mgo.Session
	db         string
	collection string
	factory    func() interface{}
}

func NewMongoReadRepository(url, database, collection string) (*MongoReadRepository, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, ErrCouldNotDialDB
	}

	session.SetMode(mgo.Strong, true)
	session.SetSafe(&mgo.Safe{W: 1})

	return NewMongoReadRepositoryWithSession(session, database, collection)
}

func NewMongoReadRepositoryWithSession(session *mgo.Session, database, collection string) (*MongoReadRepository, error) {
	if session == nil {
		return nil, ErrNoDBSession
	}

	r := &MongoReadRepository{
		session:    session,
		db:         database,
		collection: collection,
	}

	return r, nil
}

func (r *MongoReadRepository) Save(id GUID, model interface{}) error {
	sess := r.session.Copy()
	defer sess.Close()

	if _, err := sess.DB(r.db).C(r.collection).UpsertId(id, model); err != nil {
		return ErrCouldNotSaveModel
	}
	return nil
}

func (r *MongoReadRepository) Find(id GUID) (interface{}, error) {
	sess := r.session.Copy()
	defer sess.Close()

	if r.factory == nil {
		return nil, ErrModelNotSet
	}

	model := r.factory()
	err := sess.DB(r.db).C(r.collection).FindId(id).One(model)
	if err != nil {
		return nil, ErrModelNotFound
	}

	return model, nil
}

func (r *MongoReadRepository) FindCustom(callback func(*mgo.Collection) *mgo.Query) ([]interface{}, error) {
	sess := r.session.Copy()
	defer sess.Close()

	if r.factory == nil {
		return nil, ErrModelNotSet
	}

	collection := sess.DB(r.db).C(r.collection)
	query := callback(collection)

	iter := query.Iter()
	result := []interface{}{}
	model := r.factory()
	for iter.Next(model) {
		result = append(result, model)
		model = r.factory()
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *MongoReadRepository) FindAll() ([]interface{}, error) {
	sess := r.session.Copy()
	defer sess.Close()

	if r.factory == nil {
		return nil, ErrModelNotSet
	}

	iter := sess.DB(r.db).C(r.collection).Find(nil).Iter()
	result := []interface{}{}
	model := r.factory()
	for iter.Next(model) {
		result = append(result, model)
		model = r.factory()
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *MongoReadRepository) Remove(id GUID) error {
	sess := r.session.Copy()
	defer sess.Close()

	err := sess.DB(r.db).C(r.collection).RemoveId(id)
	if err != nil {
		return ErrModelNotFound
	}

	return nil
}

func (r *MongoReadRepository) SetModel(factory func() interface{}) {
	r.factory = factory
}

func (r *MongoReadRepository) SetDB(db string) {
	r.db = db
}

func (r *MongoReadRepository) Clear() error {
	if err := r.session.DB(r.db).C(r.collection).DropCollection(); err != nil {
		return ErrCouldNotClearDB
	}
	return nil
}

func (r *MongoReadRepository) Close() {
	r.session.Close()
}
