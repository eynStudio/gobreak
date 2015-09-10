package db

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/goinggo/tracelog"
	"gopkg.in/mgo.v2"
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

func NewMgoCtx() *MgoCtx {
	ctx := &MgoCtx{}
	ctx.setup()
	return ctx
}

func (this *MgoCtx) setup() {
	var err error
	cfg := loadCfg()

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

func loadCfg() (cfg *MgoCfg) {
	content, err := ioutil.ReadFile("conf/mgo.json")
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