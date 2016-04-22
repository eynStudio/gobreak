package ddd

import (
	. "github.com/eynstudio/gobreak"
)

type Cmd interface {
	ID() GUID
}

type CmdHandler interface {
	CanHandle(cmd Cmd) bool
	Handle(Cmd) error
}

type Event interface {
	ID() GUID
}
