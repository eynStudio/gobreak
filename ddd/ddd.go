package ddd


type Cmd interface {
	ID() string
	AggType() string
	CmdType() string
}

type Event interface {
	ID() string
	AggType() string
	EventType() string
}