package dddd

import (
	"log"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db"
	"github.com/eynstudio/gobreak/dddd/cmdbus"
	"github.com/eynstudio/gobreak/dddd/ddd"
	"github.com/eynstudio/gobreak/dddd/store"
	"github.com/eynstudio/gobreak/di"
)

func init() {
	log.Println("dddd init")
}

func Reg(agg ddd.Agg, repo db.Repo) {
	di.Map(repo)
	store.RegRepo(agg, repo)
	cmdbus.SetAgg(agg)
}

type IdMsg struct {
	Id GUID
}

func (p *IdMsg) ID() GUID { return p.Id }

type IdCmd struct {
	IdMsg
}

func (p *IdCmd) ID() GUID { return p.Id }

type IdEvent struct {
	Id GUID
}

func (p *IdEvent) ID() GUID { return p.Id }
