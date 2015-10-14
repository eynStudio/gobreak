package ddd

import (
	. "github.com/eynstudio/gobreak"
)

type Cmd interface {
	ID() GUID
	AggType() string
	CmdType() string
}

type Event interface {
	ID() GUID
	AggType() string
	EventType() string
}

type Repository interface {
	Load(string, GUID) (Aggregate, error)
	Save(Aggregate) error
}
