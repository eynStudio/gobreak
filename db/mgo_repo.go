package db

import (
	. "github.com/eynstudio/gobreak"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

type MgoRepo interface {
	Repo
	C() *mgo.Collection
	Page(pf *PageFilter,lst T) (pager Paging)
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

func (this *mgoRepo) All(m T) T {
	err:=this.C().Find(nil).All(m)
	if err!=nil{
		panic(err)
	}
	return m
}

func (p *mgoRepo) Page(pf *PageFilter,lst T) (pager Paging) {
	pager.Total,_=p.C().Find(nil).Count()
	p.C().Find(nil).Skip(pf.Skip()).Limit(pf.PerPage).All(lst)
	pager.Items=lst
	return
}

func (this *mgoRepo) Get(id T, m T) T {
	this.C().FindId(id).One(m)
	return m
}

func (this *mgoRepo) Save(id T, m T) {
	this.C().UpsertId(id, m)
}

func (this *mgoRepo) Del(id T) {
	this.C().RemoveId(id)
}

func Str2bson(id T) bson.ObjectId {
	if reflect.TypeOf(id).Name() == "ObjectId" {
		return id.(bson.ObjectId)
	} else if IsStrT(id) {
		return bson.ObjectIdHex(id.(string))
	} else {
		return id.(bson.ObjectId)
	}
}
