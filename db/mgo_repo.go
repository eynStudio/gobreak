package db

import (
	"github.com/eynstudio/gobreak"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoRepo struct {
	Ctx  DbCtx `di`
	name string
}

func NewMgoRepo(name string) *MgoRepo {
	return &MgoRepo{name: name}
}

func (this *MgoRepo) c() *mgo.Collection {
	return this.Ctx.GetCollection(this.name).(*mgo.Collection)
}

func (this *MgoRepo) All(m interface{}) interface{} {
	this.c().Find(nil).All(m)
	return m
}

func (this *MgoRepo) Get(id interface{}, m interface{}) interface{} {
	this.c().FindId(str2bson(id)).One(m)
	return m
}

func (this *MgoRepo) Save(id interface{}, m interface{}) {
	this.c().UpsertId(str2bson(id), m)
}

func (this *MgoRepo) Del(id interface{}) {
	this.c().RemoveId(str2bson(id))
}

func str2bson(id interface{}) bson.ObjectId {
	if gobreak.IsString(id) {
		return bson.ObjectIdHex(id.(string))
	} else {
		return id.(bson.ObjectId)
	}
}
