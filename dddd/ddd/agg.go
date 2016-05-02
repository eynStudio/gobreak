package ddd

import (
	. "github.com/eynstudio/gobreak"
)

type Agg interface {
	ID() GUID
	Root() Entity
	HandleCmd(Cmd) error
	RegistedCmds() []Cmd

	GetUncommittedEvents() []Event
	ClearUncommittedEvents()
	HasUncommittedEvents() bool
	IsDeleted() bool
}

type AggBase struct {
	uncommitted []Event
	deleted     bool
}

func (p *AggBase) HandleCmd(Cmd) error           { return nil }
func (a *AggBase) StoreEvent(event Event)        { a.uncommitted = append(a.uncommitted, event) }
func (a *AggBase) GetUncommittedEvents() []Event { return a.uncommitted }
func (a *AggBase) HasUncommittedEvents() bool    { return len(a.uncommitted) > 0 }
func (a *AggBase) ClearUncommittedEvents()       { a.uncommitted = []Event{} }
func (p *AggBase) IsDeleted() bool               { return p.deleted }
func (p *AggBase) SetDeleted()                   { p.deleted = true }
