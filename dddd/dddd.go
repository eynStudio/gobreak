package dddd

import (
	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/dddd/cmdbus"
)

func init() {
	cmdbus.SetHandler(&aggCmdHandler{})
}

type Event interface {
	ID() GUID
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
