package db

import (
	"reflect"

	. "github.com/eynstudio/gobreak"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoRepo interface {
	Repo
	C() *mgo.Collection
	Page(pf *PageFilter, lst T) (pager Paging)
	GetByQWithFields(q bson.M, fields []string, i T)
	ListByQWithFields(q bson.M, fields []string, i T)
	UpdateSetFiled(id bson.ObjectId, field string, value T)
	UpdateSetMap(id bson.ObjectId, value bson.M)
	GetWithFields(id bson.ObjectId, fields []string, i T)
}

type mgoRepo struct {
	Ctx  DbCtx `di`
	name string
}

func NewMgoRepo(name string) *mgoRepo {
	return &mgoRepo{name: name}
}

func (p *mgoRepo) C() *mgo.Collection {
	return p.Ctx.GetCollection(p.name).(*mgo.Collection)
}

func (p *mgoRepo) All(m T) T {
	err := p.C().Find(nil).All(m)
	if err != nil {
		panic(err)
	}
	return m
}

func (p *mgoRepo) Page(pf *PageFilter, lst T) (pager Paging) {
	pager.Total, _ = p.C().Find(nil).Count()
	p.C().Find(nil).Skip(pf.Skip()).Limit(pf.PerPage).All(lst)
	pager.Items = lst
	return
}

func (p *mgoRepo) Get(id T, m T) T {
	p.C().FindId(id).One(m)
	return m
}

func (p *mgoRepo) Save(id T, m T) {
	p.C().UpsertId(id, m)
}

func (p *mgoRepo) Del(id T) {
	p.C().RemoveId(id)
}

func (p *mgoRepo) GetWithFields(id bson.ObjectId, fields []string, i T) {
	p.Ctx.(*MgoCtx).GetWithFields(p.C(), id, fields, i)
}

func (p *mgoRepo) GetByQWithFields(q bson.M, fields []string, i T) {
	p.Ctx.(*MgoCtx).GetByQWithFields(p.C(), q, fields, i)
}

func (p *mgoRepo) ListByQWithFields(q bson.M, fields []string, i T) {
	p.Ctx.(*MgoCtx).ListByQWithFields(p.C(), q, fields, i)
}

func (p *mgoRepo) UpdateSetFiled(id bson.ObjectId, field string, value T) {
	p.Ctx.(*MgoCtx).UpdateSetFiled(p.C(), id, field, value)
}

func (p *mgoRepo) UpdateSetMap(id bson.ObjectId, value bson.M) {
	p.Ctx.(*MgoCtx).UpdateSetMap(p.C(), id, value)
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
