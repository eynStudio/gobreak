package mgo

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db"
	//	. "github.com/eynstudio/gobreak/ddd"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//	"reflect"
)

type MgoRepo interface {
	Repo
	NewId() GUID
	GetAs(id T, m T)
	GetQ(q T) T
	GetQAs(val T,q T)
	CopySession() *mgo.Session
	C(session *mgo.Session) *mgo.Collection
	Find(q interface{}) []T
	FindAs(lst T, q interface{})
	Has(q T) bool
	HasId(id GUID) bool
	Count(q T) int
	Page(pf *PageFilter, q interface{}) (pager Paging)
	PageAs(ptype T, pf *PageFilter, q interface{}) (pager Paging)
	GetByQWithFields(q bson.M, fields []string, i T)
	ListByQWithFields(q bson.M, fields []string, i T)
	UpdateSetFiled(id GUID, field string, value T)
	UpdateSetMap(id GUID, value bson.M)
	GetWithFields(id GUID, fields []string, i T)
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
	if session == nil {
		return p.Ctx.Db().C(p.name)
	} else {
		return p.Ctx.C(session, p.name)
	}
}

func (p *mgoRepo) All() []T {
	return p.Find(nil)
}
func (p *mgoRepo) Find(q interface{}) []T {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	iter := p.C(sess).Find(q).Iter()
	return p.fetchItems(iter)
}

func (p *mgoRepo) FindAs(lst T, q interface{}) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.C(sess).Find(q).All(&lst)
}
func (p *mgoRepo) Page(pf *PageFilter, q interface{}) (pager Paging) {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	pager.Total, _ = p.C(sess).Find(q).Count()
	iter := p.C(sess).Find(q).Skip(pf.Skip()).Limit(pf.PerPage).Iter()
	pager.Items = p.fetchItems(iter)
	return
}

func (p *mgoRepo) PageAs(pslice T, pf *PageFilter, q interface{}) (pager Paging) {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	pager.Total, _ = p.C(sess).Find(q).Count()
	p.C(sess).Find(q).Skip(pf.Skip()).Limit(pf.PerPage).All(pslice)
	pager.Items = pslice
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

func (p *mgoRepo) NewId() GUID { return NewGuid() }

func (p *mgoRepo) Get(id T) T {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	m := p.factory()
	p.C(sess).FindId(id).One(m)
	return m
}

func (p *mgoRepo) GetAs(id T, m T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	p.C(sess).FindId(id).One(m)
}

func (p *mgoRepo) GetQ(q T) T {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	m := p.factory()
	p.C(sess).Find(q).One(m)
	return m
}
func (p *mgoRepo) GetQAs(val T,q T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.C(sess).Find(q).One(val)
}

func (p *mgoRepo) Has(q T) bool {
	return p.Count(q)>0
}

func (p *mgoRepo) HasId(id GUID) bool {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	n,_:=p.C(sess).FindId(id).Count()
	return n>0
}

func (p *mgoRepo) Count(q T) int {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	n,_:=p.C(sess).Find(q).Count()
	return n
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

func (p *mgoRepo) Clear() error {
	sess := p.Ctx.CopySession()
	defer sess.Close()

	if err := p.C(sess).DropCollection(); err != nil {
		return err
	}
	return nil
}

func (p *mgoRepo) GetWithFields(id GUID, fields []string, i T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.Ctx.GetWithFields(p.C(sess), id, fields, i)
}

func (p *mgoRepo) GetByQWithFields(q bson.M, fields []string, i T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.Ctx.GetByQWithFields(p.C(sess), q, fields, i)
}

func (p *mgoRepo) ListByQWithFields(q bson.M, fields []string, i T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.Ctx.ListByQWithFields(p.C(sess), q, fields, i)
}

func (p *mgoRepo) UpdateSetFiled(id GUID, field string, value T) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.Ctx.UpdateSetFiled(p.C(sess), id, field, value)
}

func (p *mgoRepo) UpdateSetMap(id GUID, value bson.M) {
	sess := p.Ctx.CopySession()
	defer sess.Close()
	p.Ctx.UpdateSetMap(p.C(sess), id, value)
}
