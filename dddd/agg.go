package dddd

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/dddd/cmdbus"
)

type Agg interface {
	ID() GUID
	HandleCmd(Cmd) error
	RegistedCmds() []Cmd
	ApplyEvent(events Event)
}

type AggBase struct {
	Root        Entity
	uncommitted []Event
}

func (p *AggBase) ID() GUID            { return p.Root.ID() }
func (p *AggBase) HandleCmd(Cmd) error { return nil }

func (a *AggBase) StoreEvent(event Event)        { a.uncommitted = append(a.uncommitted, event) }
func (a *AggBase) GetUncommittedEvents() []Event { return a.uncommitted }
func (a *AggBase) HasUncommittedEvents() bool    { return len(a.uncommitted) > 0 }
func (a *AggBase) ClearUncommittedEvents()       { a.uncommitted = []Event{} }
