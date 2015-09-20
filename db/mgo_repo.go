package db

import (
	"github.com/eynstudio/gobreak"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	this.C().FindId(Str2bson(id)).One(m)
	return m
}

func (this *mgoRepo) Save(id interface{}, m interface{}) {
	this.C().UpsertId(Str2bson(id), m)
}

func (this *mgoRepo) Del(id interface{}) {
	this.C().RemoveId(Str2bson(id))
}

func Str2bson(id interface{}) bson.ObjectId {
	if gobreak.IsString(id) {
		return bson.ObjectIdHex(id.(string))
	} else {
		return id.(bson.ObjectId)
	}
}
