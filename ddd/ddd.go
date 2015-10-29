package ddd

import (
	. "github.com/eynstudio/gobreak"
)

type Cmd interface {
	ID() GUID
}

type Event interface {
	ID() GUID
}

type IdCmd struct {
	Id GUID
}

func (p *IdCmd) ID() GUID { return p.Id }

type IdEvent struct {
	Id GUID
}

func (p *IdEvent) ID() GUID { return p.Id }
