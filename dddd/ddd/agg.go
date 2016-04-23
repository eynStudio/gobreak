package ddd

import (
	. "github.com/eynstudio/gobreak"
)

type Agg interface {
	ID() GUID
	Root() Entity
	HandleCmd(Cmd) error
	RegistedCmds() []Cmd
	ApplyEvent(events Event) //?需要？
	GetUncommittedEvents() []Event
	ClearUncommittedEvents()
	HasUncommittedEvents() bool
}

type AggBase struct {
	uncommitted []Event
}

func (p *AggBase) HandleCmd(Cmd) error           { return nil }
func (a *AggBase) StoreEvent(event Event)        { a.uncommitted = append(a.uncommitted, event) }
func (a *AggBase) GetUncommittedEvents() []Event { return a.uncommitted }
func (a *AggBase) HasUncommittedEvents() bool    { return len(a.uncommitted) > 0 }
func (a *AggBase) ClearUncommittedEvents()       { a.uncommitted = []Event{} }
