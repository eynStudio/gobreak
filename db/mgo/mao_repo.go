package mgo

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db"
	//	. "github.com/eynstudio/gobreak/ddd"
	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
	//	"reflect"
)

type MgoRepo interface {
	Repo
	CopySession() *mgo.Session
	C(session *mgo.Session) *mgo.Collection

	Page(pf *PageFilter, lst T) (pager Paging)
	//	GetByQWithFields(q bson.M, fields []string, i T)
	//	ListByQWithFields(q bson.M, fields []string, i T)
	//	UpdateSetFiled(id bson.ObjectId, field string, value T)
	//	UpdateSetMap(id bson.ObjectId, value bson.M)
	//	GetWithFields(id bson.ObjectId, fields []string, i T)
}

type mgoRepo struct {
	Ctx     *MgoCtx `di`
	name    string
	factory func() T
}

func NewMgoRepo(name string, factory func() T) *mgoRepo {
	return &mgoRepo{name: name, factory: factory}
}

func (p *mgoRepo) CopySession() *mgo.Session {
	return p.Ctx.CopySession()
}

func (p *mgoRepo) C(session *mgo.Session) *mgo.Collection {
	return p.Ctx.C(session, p.name)
}

func (p *mgoRepo) All() []T {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	iter := p.C(sess).Find(nil).Iter()
	return p.fetchItems(iter)
}

func (p *mgoRepo) Page(pf *PageFilter, lst T) (pager Paging) {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	pager.Total, _ = p.C(sess).Find(nil).Count()
	iter := p.C(sess).Find(nil).Iter()
	pager.Items = p.fetchItems(iter)
	return
}

func (p *mgoRepo) fetchItems(iter *mgo.Iter) []T {
	result := []T{}
	model := p.factory()
	for iter.Next(model) {
		result = append(result, model)
		model = p.factory()
	}
	if err := iter.Close(); err != nil {
		return nil
	}

	return result
}
func (p *mgoRepo) Get(id T) T {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	m := p.factory()
	p.C(sess).FindId(id).One(m)
	return m
}

func (p *mgoRepo) Save(id T, m T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.C(sess).UpsertId(id, m)
}

func (p *mgoRepo) Del(id T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.C(sess).RemoveId(id)
}

//func (p *mgoRepo) GetWithFields(id bson.ObjectId, fields []string, i T) {
//	p.Ctx.(*MgoCtx).GetWithFields(p.C(), id, fields, i)
//}

//func (p *mgoRepo) GetByQWithFields(q bson.M, fields []string, i T) {
//	p.Ctx.(*MgoCtx).GetByQWithFields(p.C(), q, fields, i)
//}

//func (p *mgoRepo) ListByQWithFields(q bson.M, fields []string, i T) {
//	p.Ctx.(*MgoCtx).ListByQWithFields(p.C(), q, fields, i)
//}

//func (p *mgoRepo) UpdateSetFiled(id bson.ObjectId, field string, value T) {
//	p.Ctx.(*MgoCtx).UpdateSetFiled(p.C(), id, field, value)
//}

//func (p *mgoRepo) UpdateSetMap(id bson.ObjectId, value bson.M) {
//	p.Ctx.(*MgoCtx).UpdateSetMap(p.C(), id, value)
//}
