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
