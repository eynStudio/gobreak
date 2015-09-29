package db

import (
	"github.com/eynstudio/gobreak"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

type MgoRepo interface {
	Repo
	C() *mgo.Collection
}

type mgoRepo struct {
	Ctx  DbCtx `di`
	name string
}

func NewMgoRepo(name string) *mgoRepo {
	return &mgoRepo{name: name}
}

func (this *mgoRepo) C() *mgo.Collection {
	return this.Ctx.GetCollection(this.name).(*mgo.Collection)
}

func (this *mgoRepo) All(m interface{}) interface{} {
	this.C().Find(nil).All(m)
	return m
}

func (this *mgoRepo) Get(id interface{}, m interface{}) interface{} {
	this.C().FindId(id).One(m)
	return m
}

func (this *mgoRepo) Save(id interface{}, m interface{}) {
	this.C().UpsertId(id, m)
}

func (this *mgoRepo) Del(id interface{}) {
	this.C().RemoveId(id)
}

func Str2bson(id interface{}) bson.ObjectId {
	if reflect.TypeOf(id).Name() == "ObjectId" {
		return id.(bson.ObjectId)
	} else if gobreak.IsStrT(id) {
		return bson.ObjectIdHex(id.(string))
	} else {
		return id.(bson.ObjectId)
	}
}
