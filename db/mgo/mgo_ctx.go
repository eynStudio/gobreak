package mgo

import (
	. "github.com/eynstudio/gobreak"
	log "github.com/goinggo/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoCfg struct {
	Server string
	Db     string
	User   string
	Pwd    string
}

type MgoCtx struct {
	cfg     *MgoCfg
	session *mgo.Session
}

func NewMgoCtx(cfg *MgoCfg) *MgoCtx {
	ctx := &MgoCtx{cfg: cfg}
	var err error
	if ctx.session, err = mgo.Dial(cfg.Server); err != nil {
		log.Error(err, "Dial", "mdb.Startup")
		panic(err)
	}
	ctx.session.SetMode(mgo.Strong, true)
	ctx.session.SetSafe(&mgo.Safe{W: 1})
	return ctx
}

func (p *MgoCtx) CopySession() *mgo.Session {
	return p.session.Copy()
}

func (p *MgoCtx) C(session *mgo.Session, name string) *mgo.Collection {
	db := session.DB(p.cfg.Db)
	if err := db.Login(p.cfg.User, p.cfg.Pwd); err != nil {
		log.Error(err, "Login", "mdb.Startup")
		return nil
	}

	return db.C(name)
}

func (p *MgoCtx) GetWithFields(c *mgo.Collection, id GUID, fields []string, i T) {
	selector := Fields2BsonM(fields)
	c.FindId(id).Select(selector).One(i)
}

func (p *MgoCtx) GetByQWithFields(c *mgo.Collection, q bson.M, fields []string, i T) {
	selector := Fields2BsonM(fields)
	c.Find(q).Select(selector).One(i)
}

func (p *MgoCtx) ListByQWithFields(c *mgo.Collection, q bson.M, fields []string, i T) {
	selector := Fields2BsonM(fields)
	c.Find(q).Select(selector).All(i)
}

func (p *MgoCtx) UpdateSetFiled(c *mgo.Collection, id GUID, field string, value T) {
	p.UpdateSetMap(c, id, bson.M{field: value})
}

func (p *MgoCtx) UpdateSetMap(c *mgo.Collection, id GUID, value bson.M) {
	c.UpdateId(id, bson.M{"$set": value})
}

func Fields2BsonM(fields []string) bson.M {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	return selector
}
