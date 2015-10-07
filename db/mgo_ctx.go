package db

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/eynstudio/gobreak"
	log "github.com/goinggo/tracelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoCtx struct {
	db      *mgo.Database
	session *mgo.Session
}

type MgoCfg struct {
	Server string
	Name   string
	User   string
	Pwd    string
}

func NewMgoCtx(cfgFile string) *MgoCtx {
	if cfgFile == "" {
		cfgFile = "conf/mgo.json"
	}
	ctx := &MgoCtx{}
	ctx.setup(cfgFile)
	return ctx
}

func (this *MgoCtx) setup(cfgFile string) {
	var err error
	cfg := loadCfg(cfgFile)

	if this.session, err = mgo.Dial(cfg.Server); err != nil {
		log.Error(err, "Dial", "mdb.Startup")
		panic(err)
	}

	this.db = this.session.DB(cfg.Name)

	if err = this.db.Login(cfg.User, cfg.Pwd); err != nil {
		log.Error(err, "Login", "mdb.Startup")
		panic(err)
	}
}

func (p *MgoCtx) GetCollection(name string) T {
	return p.db.C(name)
}

func (p *MgoCtx) GetWithFields(c *mgo.Collection, id bson.ObjectId, fields []string, i T) {
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

func (p *MgoCtx) UpdateSetFiled(c *mgo.Collection, id bson.ObjectId, field string, value T) {
	p.UpdateSetMap(c, id, bson.M{field: value})
}

func (p *MgoCtx) UpdateSetMap(c *mgo.Collection, id bson.ObjectId, value bson.M) {
	c.UpdateId(id, bson.M{"$set": value})
}

func Fields2BsonM(fields []string) bson.M {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	return selector
}

func loadCfg(cfgFile string) (cfg *MgoCfg) {
	content, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(content, &cfg); err != nil {
		panic(err)
	}

	return
}

func (this *MgoCtx) Shutdown() error {
	if this.session != nil {
		this.session.Close()
	}
	return nil
}
