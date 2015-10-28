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

type IdEvent struct {
	Id GUID
}