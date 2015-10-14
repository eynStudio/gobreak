package ddd

type GUID string

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
